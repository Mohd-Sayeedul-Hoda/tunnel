package request

import (
	"context"
)

type Validator interface {
	Valid(ctx context.Context) *Valid
}

type Valid struct {
	Errors map[string]string
}

func NewValidator() *Valid {
	return &Valid{
		Errors: make(map[string]string),
	}
}

func (v *Valid) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Valid) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Valid) Valid() bool {
	return len(v.Errors) == 0
}
