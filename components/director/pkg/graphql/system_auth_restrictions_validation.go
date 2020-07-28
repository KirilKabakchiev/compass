package graphql

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

func (i SystemAuthAccessInput) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.For, validation.Required),
		validation.Field(&i.To, validation.Required),
	)
}
