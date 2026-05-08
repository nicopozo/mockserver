package model_test

import (
	"net/http"
	"testing"

	"github.com/nicopozo/mockserver/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	type args struct {
		causeCode int64
		message   string
		args      []interface{}
	}

	tests := []struct {
		name          string
		args          args
		wantStatus    int
		wantCauseDesc string
		wantMessage   string
	}{
		{
			name: "Should create internal error",
			args: args{
				causeCode: model.InternalError,
				message:   "Some error %s",
				args:      []interface{}{"detail"},
			},
			wantStatus:    http.StatusInternalServerError,
			wantCauseDesc: "Internal server error",
			wantMessage:   "Some error detail",
		},
		{
			name: "Should create validation error",
			args: args{
				causeCode: model.ValidationError,
				message:   "Invalid field",
				args:      nil,
			},
			wantStatus:    http.StatusBadRequest,
			wantCauseDesc: "Request validation failed",
			wantMessage:   "Invalid field",
		},
		{
			name: "Should create not found error",
			args: args{
				causeCode: model.ResourceNotFoundError,
				message:   "Resource not found",
				args:      nil,
			},
			wantStatus:    http.StatusNotFound,
			wantCauseDesc: "Resource Not Found",
			wantMessage:   "Resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := model.NewError(tt.args.causeCode, tt.args.message, tt.args.args...)

			assert.Equal(t, tt.wantStatus, err.Status)
			assert.Equal(t, tt.wantCauseDesc, err.ErrorCause[0].Description)
			assert.Equal(t, tt.wantMessage, err.Message)
		})
	}
}
