package model_test

import (
	"testing"

	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/stretchr/testify/assert"
)

//nolint:maintidx
func TestAssertion_Assert(t *testing.T) {
	type assertionsFields struct {
		Type  string
		Value string
		Min   float64
		Max   float64
	}

	type args struct {
		variables *model.Variable
	}

	type want struct {
		msg     string
		isValid bool
	}

	tests := []struct {
		name             string
		assertionsFields assertionsFields
		args             args
		want             want
	}{
		{
			name: "Should validate successfully type is 'present' and value is present",
			assertionsFields: assertionsFields{
				Type: "present",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "30",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'present' and value is empty",
			assertionsFields: assertionsFields{
				Type: "present",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "",
				},
			},
			want: want{
				msg:     "variable 'limit' is not present in request",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'number' and value is a valid number",
			assertionsFields: assertionsFields{
				Type: "number",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "30",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'number' and value is not a valid number",
			assertionsFields: assertionsFields{
				Type: "number",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "not_a_number",
				},
			},
			want: want{
				msg:     "variable 'limit' is not a valid number",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'string' and value is a valid string",
			assertionsFields: assertionsFields{
				Type: "string",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "username",
					Key:   "username",
					Value: "user01",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'string' and value is not a valid string",
			assertionsFields: assertionsFields{
				Type: "string",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "username",
					Key:   "username",
					Value: "01",
				},
			},
			want: want{
				msg:     "variable 'username' is not a valid string",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'equals' and value is the expected",
			assertionsFields: assertionsFields{
				Type:  "equals",
				Value: "user01",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "username",
					Key:   "username",
					Value: "user01",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'equals' and value is not the expected",
			assertionsFields: assertionsFields{
				Type:  "equals",
				Value: "user01",
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "username",
					Key:   "username",
					Value: "user02",
				},
			},
			want: want{
				msg:     "variable 'username' value is 'user02' but expected was 'user01'",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is in range",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "30",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is = Min Value",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "0",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is = Max Value",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "100",
				},
			},
			want: want{
				isValid: true,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is under Min Value",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "-100",
				},
			},
			want: want{
				msg:     "variable 'limit' is not in a valid number range",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is over Max Value",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "300",
				},
			},
			want: want{
				msg:     "variable 'limit' is not in a valid number range",
				isValid: false,
			},
		},
		{
			name: "Should validate successfully type is 'range' and value is not a valid number",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  100,
			},
			args: args{
				variables: &model.Variable{
					Type:  "query",
					Name:  "limit",
					Key:   "limit",
					Value: "not_a_number",
				},
			},
			want: want{
				msg:     "variable 'limit' is not a valid number",
				isValid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := model.Assertion{
				Type:  tt.assertionsFields.Type,
				Value: tt.assertionsFields.Value,
				Min:   tt.assertionsFields.Min,
				Max:   tt.assertionsFields.Max,
			}
			msg, isValid := a.Assert(tt.args.variables)

			assert.Equal(t, tt.want.msg, msg)
			assert.Equal(t, tt.want.isValid, isValid)
		})
	}
}

func TestAssertion_Validate(t *testing.T) {
	type assertionsFields struct {
		Type  string
		Value string
		Min   float64
		Max   float64
	}

	tests := []struct {
		name             string
		assertionsFields assertionsFields
		wantErr          error
	}{
		{
			name: "Validate successfully when type is present",
			assertionsFields: assertionsFields{
				Type: "present",
			},
			wantErr: nil,
		},
		{
			name: "Validate successfully when type is number",
			assertionsFields: assertionsFields{
				Type: "number",
			},
			wantErr: nil,
		},
		{
			name: "Validate successfully when type is string",
			assertionsFields: assertionsFields{
				Type: "string",
			},
			wantErr: nil,
		},
		{
			name: "Validate successfully when type is equals and value is set",
			assertionsFields: assertionsFields{
				Type:  "equals",
				Value: "assertion_value",
			},
			wantErr: nil,
		},
		{
			name: "Validate successfully when type is equals and value is not set",
			assertionsFields: assertionsFields{
				Type: "equals",
			},
			wantErr: mockserrors.InvalidRulesError{Message: "value is required when type is 'equals'"},
		},
		{
			name: "Validate successfully when type is range and value is correct",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  0,
				Max:  10,
			},
			wantErr: nil,
		},
		{
			name: "Validate successfully when type is range and value is not correct",
			assertionsFields: assertionsFields{
				Type: "range",
				Min:  20,
				Max:  10,
			},
			wantErr: mockserrors.InvalidRulesError{Message: "range min must be lower than max"},
		},
		{
			name: "Should return error when type is not valid",
			assertionsFields: assertionsFields{
				Type: "invalid_type",
			},
			wantErr: mockserrors.InvalidRulesError{Message: "type 'invalid_type' is not valid"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := model.Assertion{
				Type:  tt.assertionsFields.Type,
				Value: tt.assertionsFields.Value,
				Min:   tt.assertionsFields.Min,
				Max:   tt.assertionsFields.Max,
			}

			err := a.Validate()

			assert.Equal(t, tt.wantErr, err)
		})
	}
}
