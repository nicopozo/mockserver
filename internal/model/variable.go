package model

import (
	"fmt"

	mockserrors "github.com/nicopozo/mockserver/internal/errors"
)

const (
	VariableTypeBody   = "body"
	VariableTypeHeader = "header"
	VariableTypeRandom = "random"
	VariableTypeHash   = "hash"
	VariableTypeQuery  = "query"
	VariableTypePath   = "path"
)

type Variable struct {
	Type       string       `json:"type" example:"body"`
	Name       string       `json:"name" example:"nickname"`
	Key        string       `json:"key" example:"$.nickname"`
	Assertions []*Assertion `json:"assertions"`
	Value      string       `json:"-"`
}

func (variable *Variable) Validate() error {
	if variable.Type != VariableTypeBody && variable.Type != VariableTypeHeader &&
		variable.Type != VariableTypeRandom && variable.Type != VariableTypeHash &&
		variable.Type != VariableTypeQuery && variable.Type != VariableTypePath {
		return mockserrors.InvalidRulesError{
			Message: "variable Type must be 'body', 'header', 'query', 'random', 'hash', 'query' or 'path'",
		}
	}

	return validateAssertions(variable.Assertions)
}

func validateAssertions(assertions []*Assertion) error {
	if assertions == nil {
		return nil
	}

	for _, assertion := range assertions {
		if err := assertion.Validate(); err != nil {
			return fmt.Errorf("error validating assertions, %w", err)
		}
	}

	return nil
}
