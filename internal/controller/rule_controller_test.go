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
	mockserrors "github.com/nicopozo/mockserver/internal/errors"

	"github.com/nicopozo/mockserver/internal/controller"

	"github.com/golang/mock/gomock"
	"github.com/nicopozo/mockserver/internal/model"
	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
	testutils "github.com/nicopozo/mockserver/internal/utils/test"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
)

func TestRuleController_Create(t *testing.T) {
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
				Key:         "myapp_get_4016913947",
				Application: "myapp",
				Name:        "get balance",
				Path:        "/myapp/{user}/balance",
				Strategy:    "normal",
				Method:      "GET",
				Status:      "enabled",
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
			serviceCallTimes: 1, //nolint
		},
		{
			name:        "Should fail when body in invalid",
			requestFile: "../utils/test/mocks/json/create_rule_request_invalid.json",
			serviceErr:  nil,
			wantStatus:  http.StatusBadRequest,
			want:        nil,
			wantedErr: &model.Error{
				Status:  http.StatusBadRequest,
				Error:   "Bad Request",
				Message: "Invalid JSON. unexpected end of JSON input",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceCallTimes: 0, //nolint
		},
		{
			name:        "Should return 400 when service returns InvalidRulesErrorError",
			requestFile: "../utils/test/mocks/json/create_rule_request.json",
			serviceErr:  mockserrors.InvalidRulesErrorError{Message: "invalid rule"},
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
			serviceCallTimes: 1, //nolint
		},
		{
			name:        "Should return 500 when service returns unexpected error",
			requestFile: "../utils/test/mocks/json/create_rule_request.json",
			serviceErr:  errors.New("error in service"),
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
			serviceCallTimes: 1, //nolint
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, rule *model.Rule) (*model.Rule, error) { //nolint
					if tt.serviceErr != nil { //nolint
						return nil, tt.serviceErr //nolint
					}

					rule.Method = strings.ToUpper(rule.Method)
					rule.Key = fmt.Sprintf("%v_%v_%v", strings.ToLower(rule.Application),
						strings.ToLower(rule.Method), stringutils.Hash(rule.Name))

					return rule, nil
				}).Times(tt.serviceCallTimes) //nolint

			ctx, response, err := testutils.GetGinContextWithBody(tt.requestFile) //nolint
			if err != nil {
				t.Fatalf("Error reading file: %s", err.Error())
			}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock, //nolint
			}
			rc.Create(ctx)

			if tt.wantStatus != response.Status() { //nolint
				t.Errorf("Response status code is not the expected. Expected: %v - Actual: %v",
					tt.wantStatus, response.Status()) //nolint
			}

			if tt.wantedErr != nil { //nolint
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)
				if err != nil {
					t.Fatalf("Unexpected error occurred getting error from response")
				}
				if !reflect.DeepEqual(tt.wantedErr, errorResponse) { //nolint
					t.Fatalf("Error response is not the expected. Expected: %v - Actual: %v", tt.wantedErr, errorResponse) //nolint
				}
				return
			}

			rule, err := testutils.GetRuleFromResponse(response.Bytes)
			if err != nil {
				t.Fatalf("Unexpected error occurred getting rule from response")
			}

			if !reflect.DeepEqual(tt.want, rule) { //nolint
				t.Errorf("Rule response is not the expected. Expected: %v - Actual: %v", tt.want, rule) //nolint
			}
		})
	}
}

func TestRuleController_Get(t *testing.T) {
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
				Key:         "myapp_get_4016913947",
				Application: "myapp",
				Name:        "get balance",
				Path:        "/myapp/{user}/balance",
				Strategy:    "normal",
				Method:      "GET",
				Status:      "enabled",
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
			serviceCallTimes: 1, //nolint
			key:              "myapp_get_4016913947",
		},

		{
			name:       "Should return 500 when service returns error",
			serviceErr: errors.New("error in service"),
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
			serviceCallTimes: 1, //nolint
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
			serviceCallTimes: 1, //nolint
			key:              "myapp_get_4016913947",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Get(gomock.Any(), tt.key). //nolint
										DoAndReturn(func(ctx context.Context, key string) (*model.Rule, error) { //nolint
					if tt.serviceErr != nil { //nolint
						return nil, tt.serviceErr //nolint
					}

					result := &model.Rule{
						Key:         "myapp_get_4016913947",
						Application: "myapp",
						Name:        "get balance",
						Path:        "/myapp/{user}/balance",
						Strategy:    "normal",
						Method:      "GET",
						Status:      "enabled",
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
				}).Times(tt.serviceCallTimes) //nolint

			ginContext, response := testutils.GetGinContext()
			ruleKey := gin.Param{Key: "key", Value: tt.key} //nolint
			ginContext.Params = []gin.Param{ruleKey}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock, //nolint
			}
			rc.Get(ginContext)

			if tt.wantStatus != response.Status() { //nolint
				t.Errorf("Response status code is not the expected. Expected: %v - Actual: %v",
					tt.wantStatus, response.Status()) //nolint
			}

			if tt.wantedErr != nil { //nolint
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)
				if err != nil {
					t.Fatalf("Unexpected error occurred getting error from response")
				}
				if !reflect.DeepEqual(tt.wantedErr, errorResponse) { //nolint
					t.Fatalf("Error response is not the expected. Expected: %v - Actual: %v", tt.wantedErr, errorResponse) //nolint
				}
				return
			}

			rule, err := testutils.GetRuleFromResponse(response.Bytes)
			if err != nil {
				t.Fatalf("Unexpected error occurred getting rule from response")
			}

			if !reflect.DeepEqual(tt.want, rule) { //nolint
				t.Errorf("Rule response is not the expected. Expected: %v - Actual: %v", tt.want, rule) //nolint
			}
		})
	}
}

func TestRuleController_Search(t *testing.T) { //nolint
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
						Key:         "myapp_get_4016913947",
						Application: "myapp",
						Name:        "get balance",
						Path:        "/myapp/{user}/balance",
						Strategy:    "normal",
						Method:      "GET",
						Status:      "enabled",
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
						Key:         "myapp_get_123123",
						Application: "myapp",
						Name:        "get user",
						Path:        "/myapp/{user}",
						Strategy:    "normal",
						Method:      "GET",
						Status:      "enabled",
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
			serviceCallTimes: 1, //nolint
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
						Key:         "myapp_get_4016913947",
						Application: "myapp",
						Name:        "get balance",
						Path:        "/myapp/{user}/balance",
						Strategy:    "normal",
						Method:      "GET",
						Status:      "enabled",
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
						Key:         "myapp_get_123123",
						Application: "myapp",
						Name:        "get user",
						Path:        "/myapp/{user}",
						Strategy:    "normal",
						Method:      "GET",
						Status:      "enabled",
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
			serviceCallTimes: 1, //nolint
		},
		{
			name:       "Search Rules returns 400 when invalid offset",
			wantStatus: http.StatusBadRequest,
			want:       nil,
			wantedErr: &model.Error{
				Status:  http.StatusBadRequest,
				Error:   "Bad Request",
				Message: "Error parsing pagination params: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled", "offset": "invalid"},
			serviceCallTimes: 0, //nolint
		},
		{
			name:       "Search Rules returns 400 when invalid limit",
			wantStatus: http.StatusBadRequest,
			want:       nil,
			wantedErr: &model.Error{
				Status:  http.StatusBadRequest,
				Error:   "Bad Request",
				Message: "Error parsing pagination params: strconv.ParseInt: parsing \"invalid\": invalid syntax",
				ErrorCause: []model.ErrorCause{
					{
						Code:        1001,
						Description: "Request validation failed",
					},
				},
			},
			serviceErr:       nil,
			queries:          map[string]string{"status": "enabled", "limit": "invalid"},
			serviceCallTimes: 0, //nolint
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
			serviceCallTimes: 1, //nolint
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			//nolint
			ruleServiceMock.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, params map[string]interface{}, paging model.Paging) (*model.RuleList, error) { //nolint
					if tt.serviceErr != nil { //nolint
						return nil, tt.serviceErr //nolint
					}

					expectedPaging := model.Paging{Limit: tt.want.Paging.Limit, Offset: tt.want.Paging.Offset}
					if !reflect.DeepEqual(expectedPaging, paging) { //nolint
						t.Errorf("Request Paging is not the expected. Expected: %v - Actual: %v", expectedPaging, paging) //nolint
					}
					return &model.RuleList{
						Paging: model.Paging{
							Total:  2,
							Limit:  tt.want.Paging.Limit,
							Offset: tt.want.Paging.Offset,
						},
						Results: []*model.Rule{
							{
								Key:         "myapp_get_4016913947",
								Application: "myapp",
								Name:        "get balance",
								Path:        "/myapp/{user}/balance",
								Strategy:    "normal",
								Method:      "GET",
								Status:      "enabled",
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
								Key:         "myapp_get_123123",
								Application: "myapp",
								Name:        "get user",
								Path:        "/myapp/{user}",
								Strategy:    "normal",
								Method:      "GET",
								Status:      "enabled",
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
				}).Times(tt.serviceCallTimes) //nolint

			ginContext, response := testutils.GetGinContext()
			values := ginContext.Request.URL.Query()
			for key, value := range tt.queries { //nolint
				values.Add(key, value)
			}

			ginContext.Request.URL.RawQuery = values.Encode()

			rc := &controller.RuleController{
				RuleService: ruleServiceMock, //nolint
			}
			rc.Search(ginContext)

			if tt.wantStatus != response.Status() { //nolint
				t.Errorf("Response status code is not the expected. Expected: %v - Actual: %v",
					tt.wantStatus, response.Status()) //nolint
			}

			if tt.wantedErr != nil { //nolint
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)
				if err != nil {
					t.Fatalf("Unexpected error occurred getting error from response")
				}
				if !reflect.DeepEqual(tt.wantedErr, errorResponse) { //nolint
					t.Fatalf("Error response is not the expected. Expected: %v - Actual: %v", tt.wantedErr, errorResponse) //nolint
				}
				return
			}

			rules, err := testutils.GetRuleListFromResponse(response.Bytes)
			if err != nil {
				t.Fatalf("Unexpected error occurred getting ruleList from response")
			}

			if !reflect.DeepEqual(tt.want, rules) { //nolint
				t.Errorf("RuleList response is not the expected. Expected: %v - Actual: %v", tt.want, rules) //nolint
			}
		})
	}
}

func TestRuleController_Delete(t *testing.T) {
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
			serviceCallTimes: 1, //nolint
			key:              "myapp_get_4016913947",
		},

		{
			name:       "Should return 500 when service returns error",
			serviceErr: errors.New("error in service"),
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
			serviceCallTimes: 1, //nolint
			key:              "myapp_get_4016913947",
		},
		{
			name:             "Should return 204 when service returns RuleNotFoundError",
			serviceErr:       mockserrors.RuleNotFoundError{Message: "item not found in KVS"},
			wantStatus:       http.StatusNoContent,
			wantedErr:        nil,
			serviceCallTimes: 1, //nolint
			key:              "myapp_get_4016913947",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().Delete(gomock.Any(), tt.key).Return(tt.serviceErr).Times(tt.serviceCallTimes) //nolint

			ginContext, response := testutils.GetGinContext()
			ruleKey := gin.Param{Key: "key", Value: tt.key} //nolint
			ginContext.Params = []gin.Param{ruleKey}

			rc := &controller.RuleController{
				RuleService: ruleServiceMock, //nolint
			}
			rc.Delete(ginContext)

			if tt.wantStatus != response.Status() { //nolint
				t.Errorf("Response status code is not the expected. Expected: %v - Actual: %v",
					tt.wantStatus, response.Status()) //nolint
			}

			if tt.wantedErr != nil { //nolint
				errorResponse, err := testutils.GetErrorFromResponse(response.Bytes)
				if err != nil {
					t.Fatalf("Unexpected error occurred getting error from response")
				}
				if !reflect.DeepEqual(tt.wantedErr, errorResponse) { //nolint
					t.Fatalf("Error response is not the expected. Expected: %v - Actual: %v", tt.wantedErr, errorResponse) //nolint
				}
				return
			}
		})
	}
}
