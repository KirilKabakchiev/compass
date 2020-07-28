package systemauthrestrictions

import "database/sql"

type Entity struct {
	ID                    string         `db:"id"`
	SystemAuthID          string         `db:"system_auth_id"`
	ApplicationID         sql.NullString `db:"app_id"`
	RuntimeID             sql.NullString `db:"runtime_id"`
	IntegrationSystemID   sql.NullString `db:"integration_system_id"`
	ApplicationTemplateID sql.NullString `db:"application_template_id"`
}

type Collection []Entity

func (c Collection) Len() int {
	return len(c)
}
