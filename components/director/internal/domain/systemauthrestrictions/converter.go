package systemauthrestrictions

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/kyma-incubator/compass/components/director/pkg/graphql"
)

type converter struct {
}

func NewConverter() *converter {
	return &converter{}
}

func (c *converter) ToGraphQL(in *model.SystemAuthRestrictions) (*graphql.SystemAuthAccess, error) {
	if in == nil {
		return nil, nil
	}

	return &graphql.SystemAuthAccess{
		To: &graphql.SystemAuthAccessTo{
			ApplicationID:       nil,
			RuntimeID:           nil,
			IntegrationSystemID: nil,
		},
		For: &graphql.SystemAuthAccessFor{},
	}, nil
}

func (c *converter) InputFromGraphQL(in graphql.SystemAuthAccessInput) model.SystemAuthRestrictionsInput {
	return model.SystemAuthRestrictionsInput{
		To: model.SystemAuthAccessToInput{
			ApplicationID:       in.To.ApplicationID,
			RuntimeID:           in.To.RuntimeID,
			IntegrationSystemID: in.To.IntegrationSystemID,
		},
		For: model.SystemAuthAccessForInput{
			ApplicationID:         in.For.ApplicationID,
			ApplicationTemplateID: in.For.ApplicationTemplateID,
			RuntimeID:             in.For.RuntimeID,
			IntegrationSystemID:   in.For.IntegrationSystemID,
		},
	}
}

func (c *converter) ToEntity(in model.SystemAuthRestrictions) (Entity, error) {
	return Entity{
		ID:                    in.ID,
		SystemAuthID:          in.SystemAuthID,
		ApplicationID:         repo.NewNullableString(in.AppID),
		RuntimeID:             repo.NewNullableString(in.RuntimeID),
		IntegrationSystemID:   repo.NewNullableString(in.IntegrationSystemID),
		ApplicationTemplateID: repo.NewNullableString(in.AppTemplateID),
	}, nil
}

func (c *converter) FromEntity(in Entity) (model.SystemAuthRestrictions, error) {
	return model.SystemAuthRestrictions{
		ID:                  in.ID,
		SystemAuthID:        in.SystemAuthID,
		AppID:               repo.StringPtrFromNullableString(in.ApplicationID),
		RuntimeID:           repo.StringPtrFromNullableString(in.RuntimeID),
		IntegrationSystemID: repo.StringPtrFromNullableString(in.IntegrationSystemID),
		AppTemplateID:       repo.StringPtrFromNullableString(in.ApplicationTemplateID),
	}, nil
}
