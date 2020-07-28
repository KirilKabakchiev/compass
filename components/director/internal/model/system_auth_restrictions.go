package model

type SystemAuthRestrictionsInput struct {
	To  SystemAuthAccessToInput
	For SystemAuthAccessForInput
}

type SystemAuthAccessToInput struct {
	ApplicationID       *string
	RuntimeID           *string
	IntegrationSystemID *string
}

type SystemAuthAccessForInput struct {
	ApplicationID         *string
	ApplicationTemplateID *string
	RuntimeID             *string
	IntegrationSystemID   *string
}

type SystemAuthRestrictions struct {
	ID                  string
	SystemAuthID        string
	AppID               *string
	RuntimeID           *string
	IntegrationSystemID *string
	AppTemplateID       *string
}

func (i *SystemAuthRestrictionsInput) ToSystemAuthRestrictions(id, authID string) SystemAuthRestrictions {
	return SystemAuthRestrictions{
		ID:                  id,
		SystemAuthID:        authID,
		AppID:               i.For.ApplicationID,
		RuntimeID:           i.For.RuntimeID,
		IntegrationSystemID: i.For.IntegrationSystemID,
		AppTemplateID:       i.For.ApplicationTemplateID,
	}
}
