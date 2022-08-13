package service_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	mockserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen,maintidx,nosnakecase,paralleltest
func TestRuleService_Save(t *testing.T) {
	type args struct {
		rule *model.Rule
		ctx  context.Context //nolint:containedctx
	}

	tests := []struct {
		name             string
		args             args
		want             *model.Rule
		wantedErr        error
		repositoryErr    error
		serviceCallTimes int
	}{
		{
			name: "Should save rule correctly",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
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
			},
			want: &model.Rule{
				Key:      "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1,
		},
		{
			name: "Should save rule correctly with default strategy",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:  "myapp",
					Name:   "test_mock",
					Path:   "/test",
					Method: "put",
					Status: "enabled",
					Responses: []model.Response{
						{
							Body:        "{\"balance\":5000}",
							ContentType: "application/json",
							HTTPStatus:  http.StatusOK,
							Delay:       100,
						},
					},
				},
			},
			want: &model.Rule{
				Key:      "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1,
		},
		{
			name: "Should save rule correctly with default status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
					Responses: []model.Response{
						{
							Body:        "{\"balance\":5000}",
							ContentType: "application/json",
							HTTPStatus:  http.StatusOK,
							Delay:       100,
						},
					},
				},
			},
			want: &model.Rule{
				Key:      "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1,
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
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
			},
			want:             nil,
			wantedErr:        fmt.Errorf("error creating rule - %w", errors.New("error saving rule")), //nolint:goerr113
			repositoryErr:    errors.New("error saving rule"),                                         //nolint:goerr113
			serviceCallTimes: 1,
		},
		{
			name: "Should return InvalidRuleError rule is nil",
			args: args{
				ctx: mockscontext.Background(),
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesError{Message: "rule cannot be nil"},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when name is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
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
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesError{Message: "name cannot be empty"},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when path is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "",
					Strategy: "normal",
					Method:   "put",
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
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesError{Message: "path cannot be empty"},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when path does not start with '/'",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "invalid/path",
					Strategy: "normal",
					Method:   "put",
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
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesError{Message: "path must start with '/'"},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
					Status:   "invalid status",
					Responses: []model.Response{
						{
							Body:        "{\"balance\":5000}",
							ContentType: "application/json",
							HTTPStatus:  http.StatusOK,
							Delay:       100,
						},
					},
				},
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "invalid status - only 'enabled' or 'disabled' are valid values",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when invalid method",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "invalid",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "invalid is not a valid HTTP Method",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when method is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "method cannot be empty",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when strategy is invalid",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "invalid",
					Method:   "put",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "invalid rule strategy - only 'normal', 'random', 'sequential' or 'scene' are valid values",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when response is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:     "myapp",
					Name:      "test_mock",
					Path:      "/test",
					Strategy:  "normal",
					Method:    "put",
					Status:    "enabled",
					Responses: []model.Response{},
				},
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "at least one response required",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
		{
			name: "Should return InvalidRuleError rule when response has invalid http status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Group:    "myapp",
					Name:     "test_mock",
					Path:     "/test",
					Strategy: "normal",
					Method:   "put",
					Status:   "enabled",
					Responses: []model.Response{
						{
							Body:        "{\"balance\":5000}",
							ContentType: "application/json",
							HTTPStatus:  100,
							Delay:       100,
						},
					},
				},
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesError{
				Message: "100 is not a valid HTTP Status",
			},
			repositoryErr:    nil,
			serviceCallTimes: 0,
		},
	}
	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			ruleRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, rule *model.Rule) (*model.Rule, error) {
					if tt.repositoryErr != nil {
						return nil, tt.repositoryErr
					}

					rule.Key = "the_key"

					return rule, nil
				}).Times(tt.serviceCallTimes)

			ruleService, err := service.NewRuleService(ruleRepositoryMock)
			assert.Nil(t, err)

			got, err := ruleService.Save(tt.args.ctx, tt.args.rule)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err)

				return
			}

			assert.Nil(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

//nolint:funlen,nosnakecase
func TestRuleService_Get(t *testing.T) {
	type args struct {
		ctx context.Context //nolint:containedctx
		key string
	}

	tests := []struct {
		name             string
		args             args
		want             *model.Rule
		wantedErr        error
		repositoryErr    error
		repositoryRule   *model.Rule
		serviceCallTimes int
	}{
		{
			name: "Should save rule correctly",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want: &model.Rule{
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
			wantedErr:     nil,
			repositoryErr: nil,
			repositoryRule: &model.Rule{
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
			serviceCallTimes: 1,
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want:             nil,
			wantedErr:        fmt.Errorf("error getting rule, %w", errors.New("error getting rule")), //nolint:goerr113
			repositoryErr:    errors.New("error getting rule"),                                       //nolint:goerr113
			repositoryRule:   nil,
			serviceCallTimes: 1,
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			ruleRepositoryMock.EXPECT().Get(gomock.Any(), tt.args.key).
				Return(tt.repositoryRule, tt.repositoryErr).Times(tt.serviceCallTimes)

			ruleService, err := service.NewRuleService(ruleRepositoryMock)
			assert.Nil(t, err)

			got, err := ruleService.Get(tt.args.ctx, tt.args.key)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err)

				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

//nolint:funlen,nosnakecase
func TestRuleService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context //nolint:containedctx
		key string
	}

	tests := []struct {
		name             string
		args             args
		want             *model.Rule
		wantedErr        error
		repositoryErr    error
		serviceCallTimes int
	}{
		{
			name: "Should save rule correctly",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want: &model.Rule{
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1,
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want:             nil,
			wantedErr:        fmt.Errorf("error deleting rule - %w", errors.New("error deleting rule")), //nolint:goerr113
			repositoryErr:    errors.New("error deleting rule"),                                         //nolint:goerr113
			serviceCallTimes: 1,
		},
	}

	for _, tt := range tests { //nolint:paralleltest,varnamelen
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			ruleRepositoryMock.EXPECT().Delete(gomock.Any(), tt.args.key).Return(tt.repositoryErr).Times(tt.serviceCallTimes)

			ruleService, err := service.NewRuleService(ruleRepositoryMock)
			assert.Nil(t, err)

			err = ruleService.Delete(tt.args.ctx, tt.args.key)
			if tt.wantedErr != nil {
				assert.Equal(t, tt.wantedErr, err)

				return
			}

			assert.Nil(t, err)
		})
	}
}
