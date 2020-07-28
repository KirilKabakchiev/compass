package systemauthrestrictions

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
)

//TODO graphql api that allows users to grant access to systems for particular resources (apps, templates, runtimes, int_systems)

// mutation GrantSystemAccess/RevokeSystemAccess - allowed only for consumerType Users with new scope grant_system_access given to admins
//input
// to: {
// application_ids: "app2id",
// application_template_ids: "apptemplate1id",
// runtime_ids: "runtime1id",
// integration_system_ids": "another_int_system-id"
// }
// for: {
// application_ids:"app1id", "app2id",
// application_template_ids: "apptemplate1id",
// runtime_ids: "runtime1id",
// integration_system_ids": "another_int_system-id"
//}

//TODO for example a user creates an app template and grants access to int system a to operate on it
//TODO There can be a drop down in the UI to let the user creating the application pick which integration systems should be allowed to access it
// this way creating blank app, granting access and requesting OTT can be done via the same UI window
//TODO Or this can be automated by making the integration system forward the OTT it received from the LoB app to director
// this will grant the int system permissions to access the app
//TODO another idea for automation for apps created from templates would be to grant access for the app to whoever can access the template

//go:generate mockery -name=SystemAuthRestrictionsConverter -output=automock -outpkg=automock -case=underscore
type SystemAuthRestrictionsConverter interface {
	ToGraphQL(model *model.SystemAuthRestrictions) (*graphql.SystemAuthAccess, error)
	InputFromGraphQL(in graphql.SystemAuthAccessInput) model.SystemAuthRestrictionsInput
}

//go:generate mockery -name=SystemAuthRestrictionsService -output=automock -outpkg=automock -case=underscore
type SystemAuthRestrictionsService interface {
	CreateMany(ctx context.Context, in model.SystemAuthRestrictionsInput) error
	DeleteMany(ctx context.Context, in model.SystemAuthRestrictionsInput) error
}

//go:generate mockery -name=SystemAuthService -output=automock -outpkg=automock -case=underscore
type SystemAuthService interface {
	ListForObject(ctx context.Context, objectType model.SystemAuthReferenceObjectType, objectID string) ([]model.SystemAuth, error)
}

type Resolver struct {
	transact       persistence.Transactioner
	service        SystemAuthRestrictionsService
	systemAuthsSvc SystemAuthService
	converter      SystemAuthRestrictionsConverter
}

func NewResolver(transact persistence.Transactioner,
	systemAuthService SystemAuthService,
	systemAuthRestrictionsService SystemAuthRestrictionsService,
	converter SystemAuthRestrictionsConverter) *Resolver {
	return &Resolver{
		transact:       transact,
		service:        systemAuthRestrictionsService,
		systemAuthsSvc: systemAuthService,
		converter:      converter,
	}
}

func (r *Resolver) GrantSystemAccess(ctx context.Context, in graphql.SystemAuthAccessInput) (*graphql.SystemAuthAccess, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn := r.converter.InputFromGraphQL(in)

	if err := r.service.CreateMany(ctx, convertedIn); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	//TODO we might return some stuff as well
	return nil, nil
}

func (r *Resolver) RevokeSystemAccess(ctx context.Context, in graphql.SystemAuthAccessInput) (*graphql.SystemAuthAccess, error) {
	tx, err := r.transact.Begin()
	if err != nil {
		return nil, err
	}
	defer r.transact.RollbackUnlessCommitted(tx)

	ctx = persistence.SaveToContext(ctx, tx)

	convertedIn := r.converter.InputFromGraphQL(in)

	if err := r.service.DeleteMany(ctx, convertedIn); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return nil, nil
}
