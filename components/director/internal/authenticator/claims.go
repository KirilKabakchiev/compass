package authenticator

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kyma-incubator/compass/components/director/internal/consumer"
)

type Claims struct {
	Tenant         string                 `json:"tenant"`
	ExternalTenant string                 `json:"externalTenant"`
	Scopes         string                 `json:"scopes"`
	ConsumerID     string                 `json:"consumerID"`
	SystemAuthID   string                 `json:"systemAuthID"`
	ConsumerType   consumer.ConsumerType  `json:"consumerType"`
	ConsumerLevel  consumer.ConsumerLevel `json:"consumerLevel"`
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}

	return nil
}
