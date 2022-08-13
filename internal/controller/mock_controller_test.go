package controller_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nicopozo/mockserver/internal/controller"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	testutils "github.com/nicopozo/mockserver/internal/utils/test"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestMockController_Execute(t *testing.T) { //nolint:nosnakecase,paralleltest
	tests := []struct {
		name             string
		want             string
		wantStatus       int
		wantedErr        *model.Error
		serviceErr       error
		serviceResponse  *model.Response
		serviceCallTimes int
	}{
		{
			name:       "Create rule successfully",
			want:       "{\"balance\":5000}",
			wantStatus: http.StatusOK,
			wantedErr:  nil,
			serviceErr: nil,
			serviceResponse: &model.Response{
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
				Message: "No rule found for path: /test and method: get. no rule found for path",
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
			serviceResponse:  nil,
			serviceCallTimes: 1,
		},
		{
			name:       "Should return 400 when service returns InvalidRulesError",
			want:       "",
			wantStatus: http.StatusNotFound,
			wantedErr: &model.Error{
				Message: "No rule found for path: /test and method: get. no rule found for path",
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
			serviceResponse:  nil,
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
			serviceErr:       errors.New("service error"), //nolint:goerr113
			serviceResponse:  nil,
			serviceCallTimes: 1,
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockServiceMock := mocks.NewMockMockService(mockCtrl)
			defer mockCtrl.Finish()

			mockServiceMock.EXPECT().SearchResponseForRequest(gomock.Any(), gomock.Any(), "/test", gomock.Any()).
				Return(tt.serviceResponse, tt.serviceErr).Times(tt.serviceCallTimes)

			ginContext, response := testutils.GetGinContext()
			path := gin.Param{Key: "rule", Value: "/test"}
			ginContext.Params = []gin.Param{path}
			ginContext.Request.Method = "get"

			mc := &controller.MockController{
				MockService: mockServiceMock,
			}
			mc.Execute(ginContext)

			assert.Equal(t, tt.wantStatus, response.Status())

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}

			res := string(response.Bytes)

			assert.Equal(t, tt.want, res)
		})
	}
}
