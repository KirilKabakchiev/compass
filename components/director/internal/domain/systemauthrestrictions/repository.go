package systemauthrestrictions

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
)

const tableName string = `public.system_auth_restrictions`

var (
	tableColumns = []string{"id", "system_auth_id", "app_id", "runtime_id", "integration_system_id", "application_template_id"}
)

//go:generate mockery -name=Converter -output=automock -outpkg=automock -case=underscore
type EntityConverter interface {
	ToEntity(in model.SystemAuthRestrictions) (Entity, error)
	FromEntity(in Entity) (model.SystemAuthRestrictions, error)
}

type repository struct {
	creator             repo.Creator
	existsQuerierGlobal repo.ExistQuerierGlobal

	conv EntityConverter
}

func NewRepository(conv EntityConverter) *repository {
	return &repository{
		creator:             repo.NewCreator(resource.SystemAuthRestrictions, tableName, tableColumns),
		existsQuerierGlobal: repo.NewExistQuerierGlobal(resource.SystemAuthRestrictions, tableName),
		conv:                conv,
	}
}

func (r *repository) Create(ctx context.Context, item model.SystemAuthRestrictions) error {
	entity, err := r.conv.ToEntity(item)
	if err != nil {
		return errors.Wrap(err, "while converting model to entity")
	}

	return r.creator.Create(ctx, entity)
}

func (r *repository) ExistsSystemAuthRestrictionByAppIDAndAuthID(ctx context.Context, _, appID, systemAuthID string) (bool, error) {
	return r.existsQuerierGlobal.ExistsGlobal(ctx, repo.Conditions{
		repo.NewEqualCondition("system_auth_id", systemAuthID),
		repo.NewEqualCondition("app_id", appID),
	})
}

func (r *repository) ExistsSystemAuthRestrictionByRuntimeIDAndAuthID(ctx context.Context, _, runtimeID, systemAuthID string) (bool, error) {
	return r.existsQuerierGlobal.ExistsGlobal(ctx, repo.Conditions{
		repo.NewEqualCondition("system_auth_id", systemAuthID),
		repo.NewEqualCondition("runtime_id", runtimeID),
	})
}

func (r *repository) ExistsSystemAuthRestrictionByIntegrationSystemIDAndAuthID(ctx context.Context, _, integrationSystemID, systemAuthID string) (bool, error) {
	return r.existsQuerierGlobal.ExistsGlobal(ctx, repo.Conditions{
		repo.NewEqualCondition("system_auth_id", systemAuthID),
		repo.NewEqualCondition("integration_system_id", integrationSystemID),
	})
}

func (r *repository) ExistsSystemAuthRestrictionByAppTemplateIDAndAuthID(ctx context.Context, _, appTemplateID, systemAuthID string) (bool, error) {
	return r.existsQuerierGlobal.ExistsGlobal(ctx, repo.Conditions{
		repo.NewEqualCondition("system_auth_id", systemAuthID),
		repo.NewEqualCondition("application_template_id", appTemplateID),
	})
}
