package service_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
)

func TestMockService_SearchResponseForMethodAndPath(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
		path   string
	}

	tests := []struct {
		name             string
		args             args
		want             *model.Response
		wantedErr        error
		ruleServiceRule  *model.Rule
		ruleServiceErr   error
		ruleServiceTimes int
	}{
		{
			name: "Should search response successfully",
			args: args{
				ctx:    mockscontext.Background(),
				method: "PUT",
				path:   "/test",
			},
			want: &model.Response{
				Body:        "{\"balance\":5000}",
				ContentType: "application/json",
				HTTPStatus:  http.StatusOK,
				Delay:       100,
			},
			wantedErr: nil,
			ruleServiceRule: &model.Rule{
				Key:         "key123",
				Application: "myapp",
				Name:        "test_mock",
				Path:        "/test",
				Strategy:    "normal",
				Method:      "PUT",
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
			ruleServiceErr:   nil,
			ruleServiceTimes: 1,
		},
		{
			name: "Should return error when service returns error",
			args: args{
				ctx:    mockscontext.Background(),
				method: "PUT",
				path:   "/test",
			},
			want:             nil,
			wantedErr:        errors.New("error in service"),
			ruleServiceRule:  nil,
			ruleServiceErr:   errors.New("error in service"),
			ruleServiceTimes: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().SearchByMethodAndPath(tt.args.ctx, tt.args.method, tt.args.path).
				Return(tt.ruleServiceRule, tt.ruleServiceErr).Times(tt.ruleServiceTimes)

			srv := service.MockService{
				RuleService: ruleServiceMock,
			}
			got, err := srv.SearchResponseForMethodAndPath(tt.args.ctx, tt.args.method, tt.args.path)
			if (err != nil) != (tt.wantedErr != nil) {
				t.Errorf("SearchResponseForMethodAndPath() error = %v, wantedErr %v", err, tt.wantedErr != nil)
				return
			}

			if tt.wantedErr != nil {
				if !reflect.DeepEqual(tt.wantedErr, err) {
					t.Fatalf("Error is not the expected. Expected: %v - Actual: %v", tt.wantedErr, err)
				}
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Response is not the expected. Expected: %v - Actual: %v", tt.want, got)
			}
		})
	}
}
