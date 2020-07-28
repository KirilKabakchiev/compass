package consumer

import (
	"context"
	gqlgen_graphql "github.com/99designs/gqlgen/graphql"
	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name=ApplicationRepository3 -output=automock -outpkg=automock -case=underscore
type ApplicationRepository3 interface {
	GetByID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByPackageID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByDocumentID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByAPIDefinitionID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByEventDefinitionID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByWebhookID(ctx context.Context, tenantID, id string) (*model.Application, error)
	GetApplicationByPackageInstanceAuthID(ctx context.Context, tenantID, id string) (*model.Application, error)
}

type TestExtension struct {
	repo ApplicationRepository3
}

var _ interface {
	gqlgen_graphql.ResponseInterceptor
	gqlgen_graphql.OperationInterceptor
	gqlgen_graphql.FieldInterceptor
	gqlgen_graphql.HandlerExtension
} = &TestExtension{}

func (t *TestExtension) ExtensionName() string {
	return "ApproachOne"
}

func (t *TestExtension) Validate(schema gqlgen_graphql.ExecutableSchema) error {
	return nil
}

func (t *TestExtension) InterceptOperation(ctx context.Context, next gqlgen_graphql.OperationHandler) gqlgen_graphql.ResponseHandler {
	opctx := gqlgen_graphql.GetOperationContext(ctx)
	log.Printf("%d before operation %s\n", t, opctx.OperationName)
	defer func() {
		log.Printf("%d after operation %s\n", t, opctx.OperationName)
	}()
	return next(ctx)
}

func (t *TestExtension) InterceptField(ctx context.Context, next gqlgen_graphql.Resolver) (res interface{}, err error) {
	allowedConsumer, err := LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "while getting viewer from context")
	}

	if allowedConsumer.ConsumerType == User {
		return next(ctx)
	}

	fctx := gqlgen_graphql.GetFieldContext(ctx)

	log.Printf("%d entering field interceptor %s.%s with ctx consumer type %s and ctx consumer id %s\n",
		t, fctx.Object, fctx.Field.Name, allowedConsumer.ConsumerType, allowedConsumer.ConsumerID)
	defer func() {
		log.Printf("%d exiting field %s.%s with ctx consumer type %s and ctx consumer id %s\n",
			t, fctx.Object, fctx.Field.Name, allowedConsumer.ConsumerType, allowedConsumer.ConsumerID)
	}()

	var application *model.Application
	var inputID string
	tenantID, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}
	switch fctx.Field.Name {
	case "application":
		fallthrough
	case "updateApplication":
		fallthrough
	case "unregisterApplication":
		fallthrough
	case "requestOneTimeTokenForApplication":
		fallthrough
	case "requestClientCredentialsForApplication":
		if inputID == "" {
			inputID = fctx.Args["id"].(string)
		}
		fallthrough
	case "addWebhook":
		fallthrough
	case "setApplicationLabel":
		fallthrough
	case "deleteApplicationLabel":
		fallthrough
	case "addPackage":
		if inputID == "" {
			inputID = fctx.Args["applicationID"].(string)
		}

		application, err = t.repo.GetByID(ctx, tenantID, inputID)
		if err != nil {
			return nil, err
		}
	case "updatePackage":
		fallthrough
	case "deletePackage":
		fallthrough
	case "addAPIDefinitionToPackage":
		fallthrough
	case "addEventDefinitionToPackage":
		fallthrough
	case "addDocumentToPackage":
		if inputID == "" {
			inputID = fctx.Args["packageID"].(string)
		}
		if inputID == "" {
			return nil, errors.New("ID/packageID arg missing")
		}

		application, err = t.repo.GetApplicationByPackageID(ctx, tenantID, inputID)
		if err != nil {
			return nil, err
		}
	case "updateWebhook":
		fallthrough
	case "deleteWebhook":
		if inputID == "" {
			inputID = fctx.Args["webhookID"].(string)
		}
		if inputID == "" {
			return nil, errors.New("ID/packageID arg missing")
		}

		application, err = t.repo.GetApplicationByWebhookID(ctx, tenantID, inputID)
		if err != nil {
			return nil, err
		}
	}

	switch allowedConsumer.ConsumerType {
	case Application:
		if application.ID != allowedConsumer.ConsumerID {
			return nil, apperrors.NewNotFoundError(resource.Package, inputID)
		}
	case IntegrationSystem:
		//TODO commented because in latest approach column was removed
		//if application.IntegrationSystemID != nil && *application.IntegrationSystemID != allowedConsumer.ConsumerID {
		//	return nil, apperrors.NewNotFoundError(resource.Package, inputID)
		//}
	}

	return next(ctx)
}

func (t *TestExtension) InterceptResponse(ctx context.Context, next gqlgen_graphql.ResponseHandler) *gqlgen_graphql.Response {
	log.Printf("%d before response\n", t)
	defer func() {
		log.Printf("%d after response\n", t)
	}()
	return next(ctx)
}
