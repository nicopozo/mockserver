package controller_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nicopozo/mockserver/internal/controller"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
	testutils "github.com/nicopozo/mockserver/internal/utils/test"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestRuleController_Create(t *testing.T) { //nolint:nosnakecase
	t.Parallel()

	tests := []struct {
		name             string
		requestFile      string
		serviceErr       error
		want             *model.Rule
		wantStatus       int
		wantedErr        *model.Error
		serviceCallTimes int
	}{
		{
			name:        "Create rule successfully",
			requestFile: "../utils/test/mocks/json/create_rule_request.json",
			serviceErr:  nil,
			wantStatus:  http.StatusCreated,
			want: &model.Rule{
				Key:      "myapp_get_4016913947",
				Group:    "myapp",
				Name:     "get balance",
				Path:     "/myapp/{user}/balance",
				Strategy: "normal",
				Method:   "GET",
				Status:   "enabled",
				Responses: []model.Response{
					{
						Body:        "{\"balance\":5000}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       100,
					},
				},
			},
			wantedErr:        nil,
			serviceCallTimes: 1,
		},
		{
			name:        "Should fail when body is invalid",
			requestFile: "../utils/test/mocks/json/create_rule_request_invalid.json",
			serviceErr:  nil,
			wantStatus:  http.StatusBadRequest,
			want:        nil,
			wantedErr: &model.Error{
				Status:  http.StatusBadRequest,
				Error:   "Bad Request",
				Message: "Invalid JSON. error unmarshalling body, error unmarshalling reader unexpected end of JSON input",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceCallTimes: 0,
		},
		{
			name:        "Should return 400 when service returns InvalidRulesError",
			requestFile: "../utils/test/mocks/json/create_rule_request.json",
			serviceErr:  mockserrors.InvalidRulesError{Message: "invalid rule"},
			wantStatus:  http.StatusBadRequest,
			want:        nil,
			wantedErr: &model.Error{
				Status:  http.StatusBadRequest,
				Error:   "Bad Request",
				Message: "invalid rule",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceCallTimes: 1,
		},
		{
			name:        "Should return 500 when service returns unexpected error",
			requestFile: "../utils/test/mocks/json/create_rule_request.json",
			serviceErr:  errors.New("error in service"), //nolint:goerr113
			wantStatus:  http.StatusInternalServerError,
			want:        nil,
			wantedErr: &model.Error{
				Status:  http.StatusInternalServerError,
				Error:   "Internal Server Error",
				Message: "Error occurred when saving rule. error in service",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1999,
						Description: "Internal server error",
					},
				},
			},
			serviceCallTimes: 1,
		},
	}

	for _, tt := range tests { // nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
					if tt.serviceErr != nil {
						return nil, tt.serviceErr
					}

					rule.Method = strings.ToUpper(rule.Method)
					rule.Key = fmt.Sprintf("%v_%v_%v", strings.ToLower(rule.Group),
						strings.ToLower(rule.Method), stringutils.Hash(rule.Name))

					return rule, nil
				}).Times(tt.serviceCallTimes)

			ctx, response, err := testutils.GetGinContextWithBody(tt.requestFile)
			if err != nil {
				t.Fatalf("Error reading file: %s", err.Error())
			}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock,
			}
			rc.Create(ctx)

			assert.Equal(t, tt.wantStatus, response.Status())

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}

			rule, err := testutils.GetRuleFromResponse(response.Bytes)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, rule)
		})
	}
}

//nolint:funlen
func TestRuleController_Get(t *testing.T) { //nolint:nosnakecase
	t.Parallel()

	tests := []struct {
		name             string
		serviceErr       error
		want             *model.Rule
		wantStatus       int
		wantedErr        *model.Error
		serviceCallTimes int
		key              string
	}{
		{
			name:       "Get Reconcilable successfully",
			serviceErr: nil,
			want: &model.Rule{
				Key:      "myapp_get_4016913947",
				Group:    "myapp",
				Name:     "get balance",
				Path:     "/myapp/{user}/balance",
				Strategy: "normal",
				Method:   "GET",
				Status:   "enabled",
				Responses: []model.Response{
					{
						Body:        "{\"balance\":5000}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       100,
					},
				},
			},
			wantStatus:       http.StatusOK,
			wantedErr:        nil,
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},

		{
			name:       "Should return 500 when service returns error",
			serviceErr: errors.New("error in service"), //nolint:goerr113
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantedErr: &model.Error{
				Status:  http.StatusInternalServerError,
				Error:   "Internal Server Error",
				Message: "Error occurred when getting rule. error in service",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1999,
						Description: "Internal server error",
					},
				},
			},
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},
		{
			name:       "Should return 404 when service returns RuleNotFoundError",
			serviceErr: mockserrors.RuleNotFoundError{Message: "item not found in KVS"},
			want:       nil,
			wantStatus: http.StatusNotFound,
			wantedErr: &model.Error{
				Status:  http.StatusNotFound,
				Error:   "Not Found",
				Message: "item not found in KVS",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1030,
						Description: "Resource Not Found",
					},
				},
			},
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Get(gomock.Any(), tt.key).
				DoAndReturn(func(ctx context.Context, key string) (*model.Rule, error) {
					if tt.serviceErr != nil {
						return nil, tt.serviceErr
					}

					result := &model.Rule{
						Key:      "myapp_get_4016913947",
						Group:    "myapp",
						Name:     "get balance",
						Path:     "/myapp/{user}/balance",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Responses: []model.Response{
							{
								Body:        "{\"balance\":5000}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       100,
							},
						},
					}

					return result, nil
				}).Times(tt.serviceCallTimes)

			ginContext, response := testutils.GetGinContext()
			ruleKey := gin.Param{Key: "key", Value: tt.key}
			ginContext.Params = []gin.Param{ruleKey}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock,
			}
			rc.Get(ginContext)

			assert.Equal(t, tt.wantStatus, response.Status())

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}

			rule, err := testutils.GetRuleFromResponse(response.Bytes)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, rule)
		})
	}
}

//nolint
func TestRuleController_Search(t *testing.T) {
	tests := []struct {
		name             string
		wantStatus       int
		want             *model.RuleList
		wantedErr        *model.Error
		serviceErr       error
		queries          map[string]string
		serviceCallTimes int
	}{
		{
			name:       "Search Rules successfully with default paging",
			wantStatus: http.StatusOK,
			want: &model.RuleList{
				Paging: model.Paging{
					Total:  2,
					Limit:  30,
					Offset: 0,
				},
				Results: []*model.Rule{
					{
						Key:      "myapp_get_4016913947",
						Group:    "myapp",
						Name:     "get balance",
						Path:     "/myapp/{user}/balance",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Responses: []model.Response{
							{
								Body:        "{\"balance\":5000}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       100,
							},
						},
					},
					{
						Key:      "myapp_get_123123",
						Group:    "myapp",
						Name:     "get user",
						Path:     "/myapp/{user}",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Responses: []model.Response{
							{
								Body:        "{\"user\":\"nico\"}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       0,
							},
						},
					},
				},
			},
			wantedErr:        nil,
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled"},
			serviceCallTimes: 1,
		},
		{
			name:       "Search Reconcilables successfully with custom paging",
			wantStatus: http.StatusOK,
			want: &model.RuleList{
				Paging: model.Paging{
					Total:  2,
					Limit:  3,
					Offset: 1,
				},
				Results: []*model.Rule{
					{
						Key:      "myapp_get_4016913947",
						Group:    "myapp",
						Name:     "get balance",
						Path:     "/myapp/{user}/balance",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Responses: []model.Response{
							{
								Body:        "{\"balance\":5000}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       100,
							},
						},
					},
					{
						Key:      "myapp_get_123123",
						Group:    "myapp",
						Name:     "get user",
						Path:     "/myapp/{user}",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Responses: []model.Response{
							{
								Body:        "{\"user\":\"nico\"}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       0,
							},
						},
					},
				},
			},
			wantedErr:        nil,
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled", "offset": "1", "limit": "3"},
			serviceCallTimes: 1,
		},
		{
			name:       "Search Rules returns 400 when invalid offset",
			wantStatus: http.StatusBadRequest,
			want:       nil,
			wantedErr: &model.Error{
				Status: http.StatusBadRequest,
				Error:  "Bad Request",
				Message: "Error parsing pagination params: error parsing paging offset, " +
					"strconv.ParseInt: parsing \"invalid\": invalid syntax",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled", "offset": "invalid"},
			serviceCallTimes: 0,
		},
		{
			name:       "Search Rules returns 400 when invalid limit",
			wantStatus: http.StatusBadRequest,
			want:       nil,
			wantedErr: &model.Error{
				Status: http.StatusBadRequest,
				Error:  "Bad Request",
				Message: "Error parsing pagination params: error parsing paging limit, " +
					"strconv.ParseInt: parsing \"invalid\": invalid syntax",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled", "limit": "invalid"},
			serviceCallTimes: 0,
		},
		{
			name:       "Search Rules returns 500 when service returns error",
			wantStatus: http.StatusInternalServerError,
			want:       nil,
			wantedErr: &model.Error{
				Status:  http.StatusInternalServerError,
				Error:   "Internal Server Error",
				Message: "Error occurred when searching rules. error when calling service.Search()",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1999,
						Description: "Internal server error",
					},
				},
			},
			serviceErr:       errors.New("error when calling service.Search()"),
			queries:          map[string]string{"status": "enabled"},
			serviceCallTimes: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error) {
					if tt.serviceErr != nil {
						return nil, tt.serviceErr
					}

					expectedPaging := model.Paging{Limit: tt.want.Paging.Limit, Offset: tt.want.Paging.Offset}
					if !reflect.DeepEqual(expectedPaging, paging) {
						t.Errorf("Request Paging is not the expected. Expected: %v - Actual: %v", expectedPaging, paging)
					}
					return &model.RuleList{
						Paging: model.Paging{
							Total:  2,
							Limit:  tt.want.Paging.Limit,
							Offset: tt.want.Paging.Offset,
						},
						Results: []*model.Rule{
							{
								Key:      "myapp_get_4016913947",
								Group:    "myapp",
								Name:     "get balance",
								Path:     "/myapp/{user}/balance",
								Strategy: "normal",
								Method:   "GET",
								Status:   "enabled",
								Responses: []model.Response{
									{
										Body:        "{\"balance\":5000}",
										ContentType: "application/json",
										HTTPStatus:  http.StatusOK,
										Delay:       100,
									},
								},
							},
							{
								Key:      "myapp_get_123123",
								Group:    "myapp",
								Name:     "get user",
								Path:     "/myapp/{user}",
								Strategy: "normal",
								Method:   "GET",
								Status:   "enabled",
								Responses: []model.Response{
									{
										Body:        "{\"user\":\"nico\"}",
										ContentType: "application/json",
										HTTPStatus:  http.StatusOK,
										Delay:       0,
									},
								},
							},
						},
					}, nil
				}).Times(tt.serviceCallTimes)

			ginContext, response := testutils.GetGinContext()
			values := ginContext.Request.URL.Query()
			for key, value := range tt.queries {
				values.Add(key, value)
			}

			ginContext.Request.URL.RawQuery = values.Encode()

			rc := &controller.RuleController{
				RuleService: ruleServiceMock,
			}
			rc.Search(ginContext)

			assert.Equal(t, tt.wantStatus, response.Status())

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}

			rules, err := testutils.GetRuleListFromResponse(response.Bytes)

			assert.Nil(t, err)
			assert.Equal(t, tt.want, rules)
		})
	}
}

//nolint:nosnakecase,funlen
func TestRuleController_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		serviceErr       error
		wantStatus       int
		wantedErr        *model.Error
		serviceCallTimes int
		key              string
	}{
		{
			name:             "Delete Reconcilable successfully",
			serviceErr:       nil,
			wantStatus:       http.StatusNoContent,
			wantedErr:        nil,
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},

		{
			name:       "Should return 500 when service returns error",
			serviceErr: errors.New("error in service"), //nolint:goerr113
			wantStatus: http.StatusInternalServerError,
			wantedErr: &model.Error{
				Status:  http.StatusInternalServerError,
				Error:   "Internal Server Error",
				Message: "Error occurred when deleting rule. error in service",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1999,
						Description: "Internal server error",
					},
				},
			},
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},
		{
			name:             "Should return 204 when service returns RuleNotFoundError",
			serviceErr:       mockserrors.RuleNotFoundError{Message: "item not found in KVS"},
			wantStatus:       http.StatusNoContent,
			wantedErr:        nil,
			serviceCallTimes: 1,
			key:              "myapp_get_4016913947",
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Delete(gomock.Any(), tt.key).Return(tt.serviceErr).Times(tt.serviceCallTimes)

			ginContext, response := testutils.GetGinContext()
			ruleKey := gin.Param{Key: "key", Value: tt.key}
			ginContext.Params = []gin.Param{ruleKey}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock,
			}
			rc.Delete(ginContext)

			assert.Equal(t, tt.wantStatus, response.Status())

			if tt.wantedErr != nil {
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantedErr, errorResponse)

				return
			}
		})
	}
}
