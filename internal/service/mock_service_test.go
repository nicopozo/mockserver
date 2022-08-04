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

//nolint:funlen,nosnakecase
func TestMockService_SearchResponseForMethodAndPath(t *testing.T) {
	t.Parallel()

	requestMock, _ := http.NewRequest("PUT", "url", strings.NewReader("body"))

	type args struct {
		ctx     context.Context //nolint:containedctx
		request *http.Request
		path    string
		body    string
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
				ctx:     mockscontext.Background(),
				request: requestMock,
				path:    "/test",
				body:    "body",
			},
			want: &model.Response{
				Body:        "{\"balance\":5000}",
				ContentType: "application/json",
				HTTPStatus:  http.StatusOK,
				Delay:       100,
			},
			wantedErr: nil,
			ruleServiceRule: &model.Rule{
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
			ruleServiceErr:   nil,
			ruleServiceTimes: 1,
		},
		{
			name: "Should return error when service returns error",
			args: args{
				ctx:     mockscontext.Background(),
				request: requestMock,
				path:    "/test",
				body:    "body",
			},
			want:             nil,
			wantedErr:        fmt.Errorf("error searching rule, %w", errors.New("error in service")), //nolint:goerr113
			ruleServiceRule:  nil,
			ruleServiceErr:   errors.New("error in service"), //nolint:goerr113
			ruleServiceTimes: 1,
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			ruleServiceMock := mocks.NewMockIRuleService(mockCtrl)
			defer mockCtrl.Finish()

			ruleServiceMock.EXPECT().SearchByMethodAndPath(tt.args.ctx, tt.args.request.Method, tt.args.path).
				Return(tt.ruleServiceRule, tt.ruleServiceErr).Times(tt.ruleServiceTimes)

			srv, err := service.NewMockService(ruleServiceMock)
			assert.Nil(t, err)

			got, err := srv.SearchResponseForRequest(tt.args.ctx, tt.args.request, tt.args.path, tt.args.body)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err)

				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
