package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
	testutils "github.com/nicopozo/mockserver/internal/utils/test"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockService_SearchResponseForRequest(t *testing.T) {
	type args struct {
		ctx     context.Context
		request *http.Request
		path    string
		body    string
	}

	type rulesServiceCall struct {
		searchByMethodAndPathResult model.Rule
		searchByMethodAndPathErr    error
		searchByMethodAndPathTimes  int
	}

	type want struct {
		result model.Response
		err    error
	}

	tests := []struct {
		name             string
		args             []args
		rulesServiceCall []rulesServiceCall
		want             []want
	}{
		{
			name: "Should search response successfully",
			args: []args{
				{
					ctx:     mockscontext.Background(),
					request: getMockRequest(http.MethodPut, "url", "body", nil, nil),
					path:    "/test",
					body:    "body",
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: model.Rule{
						Key:      "key123",
						Group:    "myapp",
						Name:     "test_mock",
						Path:     "/test",
						Strategy: "normal",
						Method:   "PUT",
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
					searchByMethodAndPathErr:   nil,
					searchByMethodAndPathTimes: 1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"balance\":5000}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       100,
					},
					err: nil,
				},
			},
		},
		{
			name: "Should search response successfully with successive calls",
			args: []args{
				{
					ctx:     mockscontext.Background(),
					request: getMockRequest(http.MethodPut, "url", "body", nil, nil),
					path:    "/test/1",
					body:    "body",
				},
				{
					ctx:     mockscontext.Background(),
					request: getMockRequest(http.MethodPut, "url", "body", nil, nil),
					path:    "/test/2",
					body:    "body",
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: model.Rule{
						Key:      "key123",
						Group:    "myapp",
						Name:     "test_mock",
						Path:     "/test/{id}",
						Strategy: "normal",
						Method:   "GET",
						Status:   "enabled",
						Variables: []*model.Variable{
							{
								Type: "path",
								Name: "the_id",
								Key:  "id",
							},
						},
						Responses: []model.Response{
							{
								Body:        "{\"id\":{the_id}}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       100,
							},
						},
					},
					searchByMethodAndPathErr:   nil,
					searchByMethodAndPathTimes: 1,
				},
				{
					searchByMethodAndPathResult: model.Rule{
						Key:      "key123",
						Group:    "myapp",
						Name:     "test_mock",
						Path:     "/test/{id}",
						Strategy: "normal",
						Method:   "PUT",
						Status:   "enabled",
						Variables: []*model.Variable{
							{
								Type: "path",
								Name: "the_id",
								Key:  "id",
							},
						},
						Responses: []model.Response{
							{
								Body:        "{\"id\":{the_id}}",
								ContentType: "application/json",
								HTTPStatus:  http.StatusOK,
								Delay:       100,
							},
						},
					},
					searchByMethodAndPathErr:   nil,
					searchByMethodAndPathTimes: 1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"id\":1}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       100,
					},
					err: nil,
				},
				{
					result: model.Response{
						Body:        "{\"id\":2}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       100,
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when service returns error",
			args: []args{
				{
					ctx:     mockscontext.Background(),
					request: getMockRequest(http.MethodPut, "url", "body", nil, nil),
					path:    "/test",
					body:    "body",
				},
			},
			want: []want{
				{
					result: model.Response{},
					err:    fmt.Errorf("error searching rule, %w", errors.New("error in service")),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathErr:   errors.New("error in service"),
					searchByMethodAndPathTimes: 1,
				},
			},
		},
		{
			name: "Should search response and apply variables successfully",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_vars.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 554433,\"amount\": 123, \"client_id\" : \"123456789\", \"currency\" : \"ARS\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return response successfully when equals assertions passes OK",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_equals_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 554433}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when equals assertions fails and fail_on_error is enabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/555444333/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_equals_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					err: mockserrors.AssertionError{
						Errors: []string{"variable 'user_id' value is '555444333' but expected was '554433'"},
					},
				},
			},
		},
		{
			name: "Should return response successfully when equals assertions fails and fail_on_error is disabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/555444333/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_equals_assertion.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 555444333}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return response successfully when is present assertion passes OK",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_present_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 554433, \"currency\" : \"ARS\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when is present assertions fails and fail_on_error is enabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						nil,
					),
					path: "/v1/555444333/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_present_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					err: mockserrors.AssertionError{
						Errors: []string{"variable 'currency' is not present in request"},
					},
				},
			},
		},
		{
			name: "Should return response successfully when is present assertion fails and fail_on_error is disabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_present_assertion.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 554433, \"currency\" : \"ARS\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return response successfully when is number assertion passes OK",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_number_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": 554433}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when is number assertions fails and fail_on_error is enabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/not_a_number/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_number_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					err: mockserrors.AssertionError{
						Errors: []string{"variable 'user_id' is not a valid number"},
					},
				},
			},
		},
		{
			name: "Should return response successfully when is number assertions fails and fail_on_error is disabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/not_a_number/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_number_assertion.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\": \"not_a_number\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return response successfully when is string assertion passes OK",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_string_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"currency\" : \"ARS\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when is string assertion fails and fail_on_error is enabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "123",
						},
					),
					path: "/v1/334455/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_string_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					err: mockserrors.AssertionError{
						Errors: []string{"variable 'currency' is not a valid string"},
					},
				},
			},
		},
		{
			name: "Should return response successfully when is string assertions fails and fail_on_error is disabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/not_a_number/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_is_string_assertion.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"currency\" : \"ARS\"}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},

		{
			name: "Should return response successfully when range assertion passes OK",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/554433/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_range_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\" : 554433}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
		{
			name: "Should return error when is string assertion fails and fail_on_error is enabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "123",
						},
					),
					path: "/v1/99/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_range_assertion_fail_on_error.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					err: mockserrors.AssertionError{
						Errors: []string{"variable 'user_id' is not in a valid number range"},
					},
				},
			},
		},
		{
			name: "Should return response successfully when is string assertions fails and fail_on_error is disabled",
			args: []args{
				{
					ctx: mockscontext.Background(),
					request: getMockRequest(
						http.MethodPut,
						"url",
						getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
						map[string][]string{
							"Client-Id": {
								"123456789",
							},
						},
						map[string]string{
							"currency": "ARS",
						},
					),
					path: "/v1/99/payments",
					body: getJSONFromFile(t, "../utils/test/mocks/json/request_body.json"),
				},
			},
			rulesServiceCall: []rulesServiceCall{
				{
					searchByMethodAndPathResult: getRuleFromFile(t, "../utils/test/mocks/json/rule_with_range_assertion.json"),
					searchByMethodAndPathErr:    nil,
					searchByMethodAndPathTimes:  1,
				},
			},
			want: []want{
				{
					result: model.Response{
						Body:        "{\"user_id\" : 99}",
						ContentType: "application/json",
						HTTPStatus:  200,
						Delay:       0,
						Scene:       "",
					},
					err: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockRuleService(mockCtrl)
			defer mockCtrl.Finish()

			for idx := range tt.args {
				ruleServiceMock.EXPECT().SearchByMethodAndPath(tt.args[idx].ctx, tt.args[idx].request.Method, tt.args[idx].path).
					Return(tt.rulesServiceCall[idx].searchByMethodAndPathResult, tt.rulesServiceCall[idx].searchByMethodAndPathErr).
					Times(tt.rulesServiceCall[idx].searchByMethodAndPathTimes)

				srv, err := service.NewMockService(ruleServiceMock)
				assert.Nil(t, err)

				got, err := srv.SearchResponseForRequest(
					tt.args[idx].ctx, tt.args[idx].request, tt.args[idx].path, tt.args[idx].body)
				if tt.want[idx].err != nil {
					assert.Equal(t, tt.want[idx].err, err)

					return
				}

				if assert.Nil(t, err) {
					assert.Equal(t, tt.want[idx].result, got)
				}
			}
		})
	}
}

//nolint:unparam
func getMockRequest(method, url, body string, headers map[string][]string, queries map[string]string) *http.Request {
	requestMock, _ := http.NewRequest(method, url, strings.NewReader(body))

	if headers != nil {
		requestMock.Header = headers
	}

	queryValues := requestMock.URL.Query()

	for key, value := range queries {
		queryValues.Add(key, value)
	}

	requestMock.URL.RawQuery = queryValues.Encode()

	return requestMock
}

func getJSONFromFile(t *testing.T, path string) string {
	t.Helper()

	json, err := testutils.GetJSONFromFile(path)
	assert.Nil(t, err)

	return json
}

func getRuleFromFile(t *testing.T, path string) model.Rule {
	t.Helper()

	ruleStr := getJSONFromFile(t, path)
	rule := model.Rule{}

	err := jsonutils.Unmarshal(strings.NewReader(ruleStr), &rule)
	assert.Nil(t, err)

	return rule
}
