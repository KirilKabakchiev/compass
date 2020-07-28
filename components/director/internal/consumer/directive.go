package consumer

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/pkg/errors"
)

//go:generate mockery -name=SystemAuthRestrictionsRepository -output=automock -outpkg=automock -case=underscore
type SystemAuthRestrictionsRepository interface {
	ExistsSystemAuthRestrictionByAppIDAndAuthID(ctx context.Context, tenantID, id, appID string) (bool, error)
	ExistsSystemAuthRestrictionByRuntimeIDAndAuthID(ctx context.Context, tenantID, id, runtimeID string) (bool, error)
	ExistsSystemAuthRestrictionByIntegrationSystemIDAndAuthID(ctx context.Context, tenantID, id, intSysID string) (bool, error)
	ExistsSystemAuthRestrictionByAppTemplateIDAndAuthID(ctx context.Context, tenantID, id, appTemplateID string) (bool, error)
}

//go:generate mockery -name=ApplicationRepository -output=automock -outpkg=automock -case=underscore
type ApplicationRepository interface {
	ExistsApplicationByIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByPackageIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByDocumentIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByAPIDefinitionIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByEventDefinitionIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByWebhookIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
	ExistsApplicationByPackageInstanceAuthIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
}

//go:generate mockery -name=RuntimeRepository -output=automock -outpkg=automock -case=underscore
type RuntimeRepository interface {
	ExistsRuntimeByIDAndAuthID(ctx context.Context, tenantID, id, authID string) (bool, error)
}

type limitAccessDirective struct {
	providers map[string]func(ctx context.Context, tenantID, id, authID string) (bool, error)
	db        persistence.PersistenceOp
}

func NewLimitAccessDirective(db persistence.PersistenceOp, appRepo ApplicationRepository, runtimeRepo RuntimeRepository, systemAuthRestrictionsRepo SystemAuthRestrictionsRepository) *limitAccessDirective {
	return &limitAccessDirective{
		providers: map[string]func(ctx context.Context, tenantID string, id string, authID string) (bool, error){
			"GetApplicationID": func(ctx context.Context, tenantID string, id string, authID string) (bool, error) {
				consumerInfo, err := LoadFromContext(ctx)
				if err != nil {
					return false, errors.New("error missing consumer info")
				}

				// performance optimization, db call below would work too
				if consumerInfo.ConsumerType == Application {
					return consumerInfo.ConsumerID == id, nil
				}

				return appRepo.ExistsApplicationByIDAndAuthID(ctx, tenantID, id, authID)
			},
			"GetApplicationIDByPackageID":             appRepo.ExistsApplicationByPackageIDAndAuthID,
			"GetApplicationIDByDocumentID":            appRepo.ExistsApplicationByDocumentIDAndAuthID,
			"GetApplicationIDByAPIDefinitionID":       appRepo.ExistsApplicationByAPIDefinitionIDAndAuthID,
			"GetApplicationIDByEventDefinitionID":     appRepo.ExistsApplicationByEventDefinitionIDAndAuthID,
			"GetApplicationIDByWebhookID":             appRepo.ExistsApplicationByWebhookIDAndAuthID,
			"GetApplicationIDByPackageInstanceAuthID": appRepo.ExistsApplicationByPackageInstanceAuthIDAndAuthID,
			"GetApplicationIDByAuthID":                systemAuthRestrictionsRepo.ExistsSystemAuthRestrictionByAppIDAndAuthID,
			"GetRuntimeID": func(ctx context.Context, tenantID string, id string, authID string) (bool, error) {
				consumerInfo, err := LoadFromContext(ctx)
				if err != nil {
					return false, errors.New("error missing consumer info")
				}

				// performance optimization, db call below would work too
				if consumerInfo.ConsumerType == Runtime {
					return consumerInfo.ConsumerID == id, nil
				}

				return runtimeRepo.ExistsRuntimeByIDAndAuthID(ctx, tenantID, id, authID)
			},
			"GetRuntimeIDByAuthID": systemAuthRestrictionsRepo.ExistsSystemAuthRestrictionByRuntimeIDAndAuthID,
			"GetIntegrationSystemID": func(ctx context.Context, tenantID string, id string, authID string) (bool, error) {
				consumerInfo, err := LoadFromContext(ctx)
				if err != nil {
					return false, errors.New("error missing consumer info")
				}

				return consumerInfo.ConsumerType == IntegrationSystem && consumerInfo.ConsumerID == id, nil
			},
			"GetIntegrationSystemIDByAuthID":   systemAuthRestrictionsRepo.ExistsSystemAuthRestrictionByIntegrationSystemIDAndAuthID,
			"GetApplicationTemplateIDByAuthID": systemAuthRestrictionsRepo.ExistsSystemAuthRestrictionByAppTemplateIDAndAuthID,
		},
		db: db,
	}
}

func (d *limitAccessDirective) LimitAccess(ctx context.Context, _ interface{}, next graphql.Resolver, key, idField string) (interface{}, error) {
	ctx = persistence.SaveToContext(ctx, d.db)

	consumerInfo, err := LoadFromContext(ctx)
	if err != nil {
		return nil, errors.New("error missing consumer info")
	}

	if consumerInfo.ConsumerLevel == Unrestricted {
		return next(ctx)
	}

	if consumerInfo.SystemAuthID == "" {
		return nil, errors.Errorf("system auth id not found in consumer context for consumer type %s", consumerInfo.ConsumerType)
	}

	fieldContext := graphql.GetFieldContext(ctx)
	inputID := fieldContext.Args[idField].(string)

	tenantID, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "while loading tenant from context")
	}

	providerFunc, found := d.providers[key]
	if !found {
		return nil, fmt.Errorf("owner provider not found for key %s", key)
	}

	exists, err := providerFunc(ctx, tenantID, inputID, consumerInfo.SystemAuthID)
	if err != nil {
		return nil, errors.Wrapf(err, "could not provide owning entity")
	}

	if !exists {
		return nil, apperrors.NewInvalidOperationError(fmt.Sprintf("consumer of type %s with id %s is not allowed to access the requested resource",
			consumerInfo.ConsumerType, consumerInfo.ConsumerID))
	}

	return next(ctx)
}

//type ownerProvider struct {
//	applicationOwnerProviders  map[string]func(ctx context.Context, tenantID, id string) (*model.Application, error)
//	integrationSystemProviders map[string]func(ctx context.Context, tenantID, id string) (*model.IntegrationSystem, error)
//	runtimeProviders           map[string]func(ctx context.Context, tenantID, id string) (*model.Runtime, error)
//}
//
//func (op *ownerProvider) ProvideApplication(ctx context.Context, tenantID, key, id string) (*model.Application, error) {
//	ownerProviderFunc, found := op.applicationOwnerProviders[key]
//	if !found {
//		return nil, fmt.Errorf("owner provider not found for key %s", key)
//	}
//
//	return ownerProviderFunc(ctx, id, tenantID)
//}
//
//func (op *ownerProvider) ProvideIntegrationSystem(ctx context.Context, tenantID, key, id string) (*model.Application, error) {
//	ownerProviderFunc, found := op.applicationOwnerProviders[key]
//	if !found {
//		return nil, fmt.Errorf("owner provider not found for key %s", key)
//	}
//
//	return ownerProviderFunc(ctx, id, tenantID)
//}
//
//func (op *ownerProvider) ProvideRuntime(ctx context.Context, tenantID, key, id string) (*model.Application, error) {
//	ownerProviderFunc, found := op.applicationOwnerProviders[key]
//	if !found {
//		return nil, fmt.Errorf("owner provider not found for key %s", key)
//	}
//
//	return ownerProviderFunc(ctx, id, tenantID)
//}
//
//func NewOwnerProvider(repo ApplicationRepository) *ownerProvider {
//	return &ownerProvider{
//		applicationOwnerProviders: map[string]func(ctx context.Context, tenantID string, id string) (*model.Application, error){
//			"GetByID": func(ctx context.Context, tenantID string, id string) (*model.Application, error) {
//				return &model.Application{
//					ID: id,
//				}, nil
//			},
//			"GetApplicationIDByPackageID":             repo.GetApplicationByPackageID,
//			"GetApplicationIDByDocumentID":            repo.GetApplicationByDocumentID,
//			"GetApplicationIDByAPIDefinitionID":       repo.GetApplicationByAPIDefinitionID,
//			"GetApplicationIDByEventDefinitionID":     repo.GetApplicationByEventDefinitionID,
//			"GetApplicationIDByWebhookID":             repo.GetApplicationByWebhookID,
//			"GetApplicationIDByPackageInstanceAuthID": repo.GetApplicationByPackageInstanceAuthID,
//			"GetApplicationIntSystemIDByPackageID" :
//
//		},
//		integrationSystemProviders: map[string]func(ctx context.Context, tenantID string, id string) (*model.IntegrationSystem, error){
//			"GetByID": func(ctx context.Context, tenantID string, id string) (*model.IntegrationSystem, error) {
//				return &model.IntegrationSystem{
//					ID: id,
//				}, nil
//			},
//		},
//		runtimeProviders: map[string]func(ctx context.Context, tenantID string, id string) (*model.Runtime, error){
//			"GetByID": func(ctx context.Context, tenantID string, id string) (*model.Runtime, error) {
//				return &model.Runtime{
//					ID: id,
//				}, nil
//			},
//		},
//	}
//}
