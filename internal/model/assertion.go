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

func (a Assertion) Assert(variables []*Variable) (string, bool) { //nolint:cyclop
	variable := getVariable(variables, a.VariableName)
	if variable == nil {
		return fmt.Sprintf("no variable found with name '%s'", a.VariableName), false
	}

	switch a.Type {
	case AssertionTypeNumber:
		if !isNumber(strings.Trim(variable.Value, `"`)) {
			return fmt.Sprintf("variable '%s' is not a valid number", a.VariableName), false
		}
	case AssertionTypeString:
		if isNumber(variable.Value) {
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

		if v < a.Min || v > a.Max {
			return fmt.Sprintf("variable '%s' is not in a valid number range", a.VariableName), false
		}
	case AssertionTypeIsPresent:
		if strings.TrimSpace(variable.Value) == "" {
			return fmt.Sprintf("variable '%s' is not present in request", a.VariableName), false
		}
	}

	return "", true
}

func (a Assertion) Validate() error {
	switch a.Type {
	case AssertionTypeNumber, AssertionTypeString, AssertionTypeIsPresent:
		return nil
	case AssertionTypeEquals:
		if a.Value == "" {
			return mockserrors.InvalidRulesError{Message: "value is required when type is 'equals'"}
		}

		return nil
	case AssertionTypeRange:
		if a.Min >= a.Max {
			return mockserrors.InvalidRulesError{Message: "range min must be lower than max"}
		}

		return nil
	}

	return mockserrors.InvalidRulesError{Message: fmt.Sprintf("type '%s' is not valid", a.Type)}
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

func isNumber(value string) bool {
	if _, err := strconv.ParseFloat(value, floatSize); err != nil {
		return false
	}

	return true
}
