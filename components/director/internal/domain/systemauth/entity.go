package systemauth

import "database/sql"

type Entity struct {
	ID                  string         `db:"id"`
	TenantID            sql.NullString `db:"tenant_id"`
	AppID               sql.NullString `db:"app_id"`
	RuntimeID           sql.NullString `db:"runtime_id"`
	IntegrationSystemID sql.NullString `db:"integration_system_id"`
	Value               sql.NullString `db:"value"`
	AccessLevel         string         `db:"access_level"`
}

type Collection []Entity

func (c Collection) Len() int {
	return len(c)
}
