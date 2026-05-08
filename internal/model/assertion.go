package model

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/xeipuuv/gojsonschema"
)

const (
	AssertionTypeEquals     = "equals"
	AssertionTypeNotEquals  = "not_equals"
	AssertionTypeString     = "string"
	AssertionTypeNumber     = "number"
	AssertionTypeRange      = "range"
	AssertionTypeIsPresent  = "present"
	AssertionTypeRegex      = "regex"
	AssertionTypeContains   = "contains"
	AssertionTypeStartsWith = "starts_with"
	AssertionTypeEndsWith   = "ends_with"
	AssertionTypeLength     = "length"
	AssertionTypeIsOneOf    = "enum"
	AssertionTypeIsBoolean  = "boolean"
	AssertionTypeJSONSchema = "json_schema"

	floatSize = 64
)

type Assertion struct {
	FailOnError bool    `json:"fail_on_error"`
	Type        string  `json:"type"`
	Value       string  `json:"value"`
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
}

func (a Assertion) Assert(variable *Variable) (string, bool) { //nolint:cyclop,funlen,gocognit,gocyclo
	val := variable.Value

	switch a.Type {
	case AssertionTypeNumber:
		if !isNumber(val) {
			return fmt.Sprintf("variable '%s' is not a valid number", variable.Name), false
		}
	case AssertionTypeString:
		if isNumber(variable.Value) {
			return fmt.Sprintf("variable '%s' is not a valid string", variable.Name), false
		}
	case AssertionTypeEquals:
		if strings.TrimSpace(val) != strings.TrimSpace(a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but expected was '%s'",
				variable.Name, variable.Value, a.Value), false
		}
	case AssertionTypeNotEquals:
		if strings.TrimSpace(val) == strings.TrimSpace(a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' and it is not expected to be equal to '%s'",
				variable.Name, variable.Value, a.Value), false
		}
	case AssertionTypeRange:
		v, err := strconv.ParseFloat(val, floatSize)
		if err != nil {
			return fmt.Sprintf("variable '%s' is not a valid number", variable.Name), false
		}

		if v < a.Min || v > a.Max {
			return fmt.Sprintf("variable '%s' is not in a valid number range", variable.Name), false
		}
	case AssertionTypeIsPresent:
		if strings.TrimSpace(val) == "" {
			return fmt.Sprintf("variable '%s' is not present in request", variable.Name), false
		}
	case AssertionTypeRegex:
		matched, err := regexp.MatchString(a.Value, val)
		if err != nil {
			return fmt.Sprintf("invalid regex pattern: %s", a.Value), false
		}

		if !matched {
			return fmt.Sprintf("variable '%s' value is '%s' but does not match regex '%s'",
				variable.Name, val, a.Value), false
		}
	case AssertionTypeContains:
		if !strings.Contains(val, a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but does not contain '%s'",
				variable.Name, val, a.Value), false
		}
	case AssertionTypeStartsWith:
		if !strings.HasPrefix(val, a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but does not start with '%s'",
				variable.Name, val, a.Value), false
		}
	case AssertionTypeEndsWith:
		if !strings.HasSuffix(val, a.Value) {
			return fmt.Sprintf("variable '%s' value is '%s' but does not end with '%s'",
				variable.Name, val, a.Value), false
		}
	case AssertionTypeLength:
		if float64(len(val)) < a.Min || float64(len(val)) > a.Max {
			return fmt.Sprintf("variable '%s' length is %d but expected range is [%.0f, %.0f]",
				variable.Name, len(val), a.Min, a.Max), false
		}
	case AssertionTypeIsOneOf:
		options := strings.Split(a.Value, ",")
		found := false

		for _, opt := range options {
			if strings.TrimSpace(opt) == strings.TrimSpace(val) {
				found = true

				break
			}
		}

		if !found {
			return fmt.Sprintf("variable '%s' value is '%s' but must be one of [%s]",
				variable.Name, val, a.Value), false
		}
	case AssertionTypeIsBoolean:
		lowered := strings.ToLower(strings.TrimSpace(val))
		if lowered != "true" && lowered != "false" {
			return fmt.Sprintf("variable '%s' value is '%s' but must be a boolean",
				variable.Name, val), false
		}
	case AssertionTypeJSONSchema:
		schemaLoader := gojsonschema.NewStringLoader(a.Value)
		documentLoader := gojsonschema.NewStringLoader(val)

		result, err := gojsonschema.Validate(schemaLoader, documentLoader)
		if err != nil {
			return fmt.Sprintf("invalid JSON schema or document: %v", err), false
		}

		if !result.Valid() {
			var errMsgs []string
			for _, desc := range result.Errors() {
				errMsgs = append(errMsgs, desc.String())
			}

			return fmt.Sprintf("variable '%s' does not match JSON schema: %s",
				variable.Name, strings.Join(errMsgs, "; ")), false
		}
	}

	return "", true
}

func (a Assertion) Validate() error { //nolint:cyclop
	switch a.Type {
	case AssertionTypeNumber, AssertionTypeString, AssertionTypeIsPresent:
		return nil
	case AssertionTypeEquals, AssertionTypeNotEquals, AssertionTypeRegex, AssertionTypeContains,
		AssertionTypeStartsWith, AssertionTypeEndsWith, AssertionTypeIsOneOf, AssertionTypeJSONSchema:
		if a.Value == "" {
			return mockserrors.InvalidRulesError{Message: fmt.Sprintf("value is required when type is '%s'", a.Type)}
		}

		if a.Type == AssertionTypeRegex {
			if _, err := regexp.Compile(a.Value); err != nil {
				return mockserrors.InvalidRulesError{Message: fmt.Sprintf("invalid regex pattern: %s", a.Value)}
			}
		}

		if a.Type == AssertionTypeJSONSchema {
			schemaLoader := gojsonschema.NewStringLoader(a.Value)
			if _, err := schemaLoader.LoadJSON(); err != nil {
				return mockserrors.InvalidRulesError{Message: fmt.Sprintf("invalid JSON schema: %v", err)}
			}
		}

		return nil
	case AssertionTypeLength:
		if a.Min > a.Max {
			return mockserrors.InvalidRulesError{Message: "length min must be lower than or equal to max"}
		}

		return nil
	case AssertionTypeIsBoolean:
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
	AssertionErrors []string
}

func (e *AssertionResult) AddAssertionError(failsOnError bool, assertionError string) {
	e.Fail = e.Fail || failsOnError
	e.AssertionErrors = append(e.AssertionErrors, assertionError)
}

func (e *AssertionResult) GetError() error {
	return mockserrors.AssertionError{
		Errors: e.AssertionErrors,
	}
}

func (e *AssertionResult) Print(ctx context.Context) {
	logger := mockscontext.Logger(ctx)

	for _, err := range e.AssertionErrors {
		logger.Info(e, nil, "variable is not the expected: %s", err)
	}
}

func isNumber(value string) bool {
	if _, err := strconv.ParseFloat(value, floatSize); err != nil {
		return false
	}

	return true
}
