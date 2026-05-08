package controller_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/nicopozo/mockserver/internal/controller"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	testutils "github.com/nicopozo/mockserver/internal/utils/test"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMockController_Execute(t *testing.T) {
	tests := []struct {
		name             string
		want             string
		wantStatus       int
		wantedErr        *model.Error
		serviceErr       error
		serviceResponse  model.Response
		serviceCallTimes int
		rulePath         string
	}{
		{
			name:       "Create rule successfully",
			want:       "{\"balance\":5000}",
			wantStatus: http.StatusOK,
			wantedErr:  nil,
			serviceErr: nil,
			serviceResponse: model.Response{
				Body:        "{\"balance\":5000}",
				ContentType: "application/json",
				HTTPStatus:  http.StatusOK,
				Delay:       0,
			},
			serviceCallTimes: 1,
		},
		{
			name:       "Should return 404 when no response found",
			want:       "",
			wantStatus: http.StatusNotFound,
			wantedErr: &model.Error{
				Message: "No rule found for path: /test and method: GET. no rule found for path",
				Error:   "Not Found",
				Status:  http.StatusNotFound,
				ErrorCause: []model.ErrorCause{
					{
						Code:        1030,
						Description: "Resource Not Found",
					},
				},
			},
			serviceErr: mockserrors.RuleNotFoundError{
				Message: "no rule found for path",
			},
			serviceResponse:  model.Response{},
			serviceCallTimes: 1,
		},
		{
			name:       "Should return 400 when service returns InvalidRulesError",
			want:       "",
			wantStatus: http.StatusNotFound,
			wantedErr: &model.Error{
				Message: "invalid rule",
				Error:   "Bad Request",
				Status:  http.StatusBadRequest,
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr: mockserrors.InvalidRulesError{
				Message: "invalid rule",
			},
			serviceResponse:  model.Response{},
			serviceCallTimes: 1,
		},
		{
			name:       "Should return 500 when service returns error",
			want:       "",
			wantStatus: http.StatusInternalServerError,
			wantedErr: &model.Error{
				Message: "Error occurred when getting rule. service error",
				Error:   "Internal Server Error",
				Status:  http.StatusInternalServerError,
				ErrorCause: []model.ErrorCause{
					{
						Code:        1999,
						Description: "Internal server error",
					},
				},
			},
			serviceErr:       errors.New("service error"),
			serviceResponse:  model.Response{},
			serviceCallTimes: 1,
		},
		{
			name:       "Should return 400 when service returns InvalidRulesError",
			want:       "",
			wantStatus: http.StatusBadRequest,
			wantedErr: &model.Error{
				Message: "[\"error description\"]",
				Error:   "Bad Request",
				Status:  http.StatusBadRequest,
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr: mockserrors.AssertionError{
				Errors: []string{"error description"},
			},
			serviceResponse:  model.Response{},
			serviceCallTimes: 1,
		},
		{
			name:       "Should prepend leading slash to path if missing",
			want:       "{\"balance\":5000}",
			wantStatus: http.StatusOK,
			serviceResponse: model.Response{
				Body:        "{\"balance\":5000}",
				ContentType: "application/json",
				HTTPStatus:  http.StatusOK,
			},
			serviceCallTimes: 1,
			rulePath:         "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)

			mockServiceMock := mocks.NewMockMockService(mockCtrl)
			defer mockCtrl.Finish()

			expectedPath := tt.rulePath
			if expectedPath == "" {
				expectedPath = "/test"
			}

			// Ensure expectedPath starts with / for the mock expectation
			mockSearchPath := expectedPath
			if !strings.HasPrefix(mockSearchPath, "/") {
				mockSearchPath = "/" + mockSearchPath
			}

			mockServiceMock.EXPECT().SearchResponseForRequest(gomock.Any(), gomock.Any(), mockSearchPath, gomock.Any()).
				Return(tt.serviceResponse, model.AssertionResult{}, tt.serviceErr).Times(tt.serviceCallTimes)

			response, request := testutils.GetHTTPContext()
			request.SetPathValue("rule", expectedPath)
			request.Method = http.MethodGet

			mc := &controller.MockController{
				MockService: mockServiceMock,
			}
			mc.Execute(response, request)

			assert.Equal(t, tt.wantStatus, response.Code)

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Body.Bytes())

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}

			res := response.Body.String()

			assert.Equal(t, tt.want, res)
		})
	}
}
