package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
)

const (
	AssertionTypeEquals    = "equals"
	AssertionTypeString    = "string"
	AssertionTypeNumber    = "number"
	AssertionTypeRange     = "range"
	AssertionTypeIsPresent = "present"

	floatSize = 64
)

type Assertion struct {
	FailOnError  bool    `json:"fail_on_error"`
	VariableName string  `json:"variable_name"`
	Type         string  `json:"type"`
	Value        string  `json:"value"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
}

func (a Assertion) IsValid(variables []*Variable) (string, bool) {
	variable := getVariable(variables, a.VariableName)
	if variable == nil {
		return fmt.Sprintf("no variable found with name '%s'", a.VariableName), false
	}

	switch a.Type {
	case AssertionTypeNumber:
		if _, err := strconv.ParseFloat(strings.Trim(variable.Value, `"`), floatSize); err != nil {
			return fmt.Sprintf("variable '%s' is not a valid number", a.VariableName), false
		}
	case AssertionTypeString:
		if _, err := strconv.ParseFloat(variable.Value, floatSize); err == nil {
			return fmt.Sprintf("variable '%s' is not a valid string", a.VariableName), false
		}
	case AssertionTypeEquals:
		if strings.TrimSpace(variable.Value) != strings.TrimSpace(a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but expected was '%s'",
				a.VariableName, variable.Value, a.Value), false
		}
	case AssertionTypeRange:
		v, err := strconv.ParseFloat(variable.Value, floatSize)
		if err != nil {
			return fmt.Sprintf("variable '%s' is not a valid number", a.VariableName), false
		}

		if v < a.Min || v > a.Min {
			return fmt.Sprintf("variable '%s' is not in a valid number range", a.VariableName), false
		}
	case AssertionTypeIsPresent:
		if strings.TrimSpace(variable.Value) == "" {
			return fmt.Sprintf("variable '%s' is not present in request", a.VariableName), false
		}
	}

	return "", true
}

type AssertionResult struct {
	Fail            bool
	assertionErrors []string
}

func (e *AssertionResult) AddAssertionError(failsOnError bool, assertionError string) {
	e.Fail = e.Fail || failsOnError
	e.assertionErrors = append(e.assertionErrors, assertionError)
}

func (e *AssertionResult) GetError() error {
	return mockserrors.AssertionError{
		Errors: e.assertionErrors,
	}
}

func (e *AssertionResult) Print(ctx context.Context) {
	logger := mockscontext.Logger(ctx)

	for _, err := range e.assertionErrors {
		logger.Info(e, nil, "variable is not the expected: %s", err)
	}

}

func getVariable(variables []*Variable, name string) *Variable {
	for _, variable := range variables {
		if variable.Name == name {
			return variable
		}
	}

	return nil
}
