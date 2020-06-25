package model

import mockserrors "github.com/nicopozo/mockserver/internal/errors"

const (
	VariableTypeBody   = "body"
	VariableTypeHeader = "header"
	VariableTypeRandom = "random"
	VariableTypeHash   = "hash"
	VariableTypeQuery  = "query"
)

type Variable struct {
	Type string `json:"type" example:"body"`
	Name string `json:"name" example:"nickname"`
	Key  string `json:"key" example:"$.nickname"`
}

func (variable *Variable) Validate() error {
	if variable.Type != VariableTypeBody && variable.Type != VariableTypeHeader &&
		variable.Type != VariableTypeRandom && variable.Type != VariableTypeHash &&
		variable.Type != VariableTypeQuery {
		return mockserrors.InvalidRulesError{
			Message: "variable Type must be 'body', 'header', 'query', 'random' or 'hash'",
		}
	}

	return nil
}
