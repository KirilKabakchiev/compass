package tenantmapping

import (
	"github.com/kyma-incubator/compass/components/director/internal/consumer"
)

const (
	clientCredentialScopesPrefix = "clientCredentialsRegistrationScopes"
)

type TenantContext struct {
	ExternalTenantID string
	TenantID         string
}

func NewTenantContext(externalTenantID, tenantID string) TenantContext {
	return TenantContext{
		ExternalTenantID: externalTenantID,
		TenantID:         tenantID,
	}
}

type ObjectContext struct {
	TenantContext
	Scopes        string
	SystemAuthID  string
	ConsumerID    string
	ConsumerType  consumer.ConsumerType
	ConsumerLevel consumer.ConsumerLevel
}

func NewObjectContext(tenantCtx TenantContext, scopes, consumerID string, systemAuthID string, consumerType consumer.ConsumerType, consumerLevel consumer.ConsumerLevel) ObjectContext {
	return ObjectContext{
		TenantContext: tenantCtx,
		Scopes:        scopes,
		SystemAuthID:  systemAuthID,
		ConsumerID:    consumerID,
		ConsumerType:  consumerType,
		ConsumerLevel: consumerLevel,
	}
}
