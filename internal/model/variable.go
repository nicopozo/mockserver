package model

import (
	"fmt"

	mockserrors "github.com/nicopozo/mockserver/internal/errors"
)

const (
	VariableTypeBody          = "body"
	VariableTypeXML           = "xml"
	VariableTypeHeader        = "header"
	VariableTypeRandom        = "random"
	VariableTypeRandomInt     = "random_int"
	VariableTypeRandomDecimal = "random_decimal"
	VariableTypeHash          = "hash"
	VariableTypeQuery         = "query"
	VariableTypePath          = "path"
)

type Variable struct {
	Type       string       `json:"type" example:"body"`
	Name       string       `json:"name" example:"nickname"`
	Key        string       `json:"key" example:"$.nickname"`
	Min        *float64     `json:"min,omitempty"`
	Max        *float64     `json:"max,omitempty"`
	Decimals   *int         `json:"decimals,omitempty"`
	Assertions []*Assertion `json:"assertions"`
	Value      string       `json:"-"`
}

func (variable *Variable) Validate() error {
	if variable.Type != VariableTypeBody && variable.Type != VariableTypeXML &&
		variable.Type != VariableTypeHeader && variable.Type != VariableTypeRandom &&
		variable.Type != VariableTypeHash && variable.Type != VariableTypeQuery &&
		variable.Type != VariableTypePath && variable.Type != VariableTypeRandomInt &&
		variable.Type != VariableTypeRandomDecimal {
		return mockserrors.InvalidRulesError{
			Message: "variable Type must be 'body', 'xml', 'header', 'query', 'random', " +
				"'hash', 'path', 'random_int' or 'random_decimal'",
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
