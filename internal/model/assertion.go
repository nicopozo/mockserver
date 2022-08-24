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
	FailOnError bool    `json:"fail_on_error"`
	Type        string  `json:"type"`
	Value       string  `json:"value"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
}

func (a Assertion) Assert(variable *Variable) (string, bool) { //nolint:cyclop
	switch a.Type {
	case AssertionTypeNumber:
		if !isNumber(strings.Trim(variable.Value, `"`)) {
			return fmt.Sprintf("variable '%s' is not a valid number", variable.Name), false
		}
	case AssertionTypeString:
		if isNumber(variable.Value) {
			return fmt.Sprintf("variable '%s' is not a valid string", variable.Name), false
		}
	case AssertionTypeEquals:
		if strings.TrimSpace(variable.Value) != strings.TrimSpace(a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but expected was '%s'",
				variable.Name, variable.Value, a.Value), false
		}
	case AssertionTypeRange:
		v, err := strconv.ParseFloat(variable.Value, floatSize)
		if err != nil {
			return fmt.Sprintf("variable '%s' is not a valid number", variable.Name), false
		}

		if v < a.Min || v > a.Max {
			return fmt.Sprintf("variable '%s' is not in a valid number range", variable.Name), false
		}
	case AssertionTypeIsPresent:
		if strings.TrimSpace(variable.Value) == "" {
			return fmt.Sprintf("variable '%s' is not present in request", variable.Name), false
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

func isNumber(value string) bool {
	if _, err := strconv.ParseFloat(value, floatSize); err != nil {
		return false
	}

	return true
}
