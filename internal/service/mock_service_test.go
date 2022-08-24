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
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockService_SearchResponseForMethodAndPath(t *testing.T) {
	requestMock, _ := http.NewRequest("PUT", "url", strings.NewReader("body"))

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
					request: requestMock,
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
					request: requestMock,
					path:    "/test/1",
					body:    "body",
				},
				{
					ctx:     mockscontext.Background(),
					request: requestMock,
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
					request: requestMock,
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

				assert.Nil(t, err)
				assert.Equal(t, tt.want[idx].result, got)
			}
		})
	}
}
