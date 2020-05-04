package service_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	mockserrors "github.com/nicopozo/mockserver/internal/errors"

	"github.com/golang/mock/gomock"
	mockscontext "github.com/nicopozo/mockserver/internal/context"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/service"
	"github.com/nicopozo/mockserver/internal/utils/test/mocks"
)

func TestRuleService_Save(t *testing.T) {
	type args struct {
		rule *model.Rule
		ctx  context.Context
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
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
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
			},
			want: &model.Rule{
				Key:         "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should save rule correctly with default strategy",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Method:      "put",
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
			},
			want: &model.Rule{
				Key:         "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should save rule correctly with default status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
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
				Key:         "the_key",
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
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
			},
			want:             nil,
			wantedErr:        errors.New("error saving rule"),
			repositoryErr:    errors.New("error saving rule"),
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should return InvalidRuleError rule is nil",
			args: args{
				ctx: mockscontext.Background(),
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesErrorError{Message: "rule cannot be nil"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when name is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
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
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesErrorError{Message: "name cannot be empty"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when path is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "",
					Strategy:    "normal",
					Method:      "put",
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
			},
			want:             nil,
			wantedErr:        mockserrors.InvalidRulesErrorError{Message: "path cannot be empty"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
					Status:      "invalid status",
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
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "invalid status - only 'enabled' or 'disabled' are valid values"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when invalid method",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "invalid",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "invalid is not a valid HTTP Method"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when method is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "method cannot be empty"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when method is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "invalid",
					Method:      "put",
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
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "invalid rule strategy - only 'normal', 'random' or 'sequential' are valid values"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when response is empty",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
					Status:      "enabled",
					Responses:   []model.Response{},
				},
			},
			want: nil,
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "at least one response required"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
		{
			name: "Should return InvalidRuleError rule when response has invalid http status",
			args: args{
				ctx: mockscontext.Background(),
				rule: &model.Rule{
					Application: "myapp",
					Name:        "test_mock",
					Path:        "/test",
					Strategy:    "normal",
					Method:      "put",
					Status:      "enabled",
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
			wantedErr: mockserrors.InvalidRulesErrorError{
				Message: "100 is not a valid HTTP Status"},
			repositoryErr:    nil,
			serviceCallTimes: 0, //nolint
		},
	}

	for _, tt := range tests { //nolint
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockIRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			ruleRepositoryMock.EXPECT().Save(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, rule *model.Rule) (*model.Rule, error) { //nolint
					if tt.repositoryErr != nil { //nolint
						return nil, tt.repositoryErr //nolint
					}

					rule.Key = "the_key"

					return rule, nil //nolint
				}).Times(tt.serviceCallTimes) //nolint

			ruleService := &service.RuleService{RuleRepository: ruleRepositoryMock}

			got, err := ruleService.Save(tt.args.ctx, tt.args.rule) //nolint
			if (err != nil) != (tt.wantedErr != nil) {              //nolint
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantedErr != nil) //nolint
				return
			}

			if tt.wantedErr != nil { //nolint
				if !reflect.DeepEqual(tt.wantedErr, err) { //nolint
					t.Fatalf("Error is not the expected. Expected: %v - Actual: %v", tt.wantedErr, err) //nolint
				}
				return
			}

			if !reflect.DeepEqual(tt.want, got) { //nolint
				t.Errorf("Rule is not the expected. Expected: %v - Actual: %v", tt.want, got) //nolint
			}
		})
	}
}

func TestRuleService_Get(t *testing.T) {
	type args struct {
		ctx context.Context
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
			wantedErr:     nil,
			repositoryErr: nil,
			repositoryRule: &model.Rule{
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
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want:             nil,
			wantedErr:        errors.New("error getting rule"),
			repositoryErr:    errors.New("error getting rule"),
			repositoryRule:   nil,
			serviceCallTimes: 1, //nolint
		},
	}

	for _, tt := range tests { //nolint
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockIRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			//nolint
			ruleRepositoryMock.EXPECT().Get(gomock.Any(), tt.args.key).
				Return(tt.repositoryRule, tt.repositoryErr).Times(tt.serviceCallTimes) //nolint

			ruleService := &service.RuleService{RuleRepository: ruleRepositoryMock}

			got, err := ruleService.Get(tt.args.ctx, tt.args.key) //nolint
			if (err != nil) != (tt.wantedErr != nil) {            //nolint
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantedErr != nil) //nolint
				return
			}

			if tt.wantedErr != nil { //nolint
				if !reflect.DeepEqual(tt.wantedErr, err) { //nolint
					t.Fatalf("Error is not the expected. Expected: %v - Actual: %v", tt.wantedErr, err) //nolint
				}
				return
			}

			if !reflect.DeepEqual(tt.want, got) { //nolint
				t.Errorf("Rule is not the expected. Expected: %v - Actual: %v", tt.want, got) //nolint
			}
		})
	}
}

func TestRuleService_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
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
			wantedErr:        nil,
			repositoryErr:    nil,
			serviceCallTimes: 1, //nolint
		},
		{
			name: "Should return error when repository returns error",
			args: args{
				ctx: mockscontext.Background(),
				key: "key123",
			},
			want:             nil,
			wantedErr:        errors.New("error deleting rule"),
			repositoryErr:    errors.New("error deleting rule"),
			serviceCallTimes: 1, //nolint
		},
	}

	for _, tt := range tests { //nolint
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			ruleRepositoryMock := mocks.NewMockIRuleRepository(mockCtrl)
			defer mockCtrl.Finish()

			ruleRepositoryMock.EXPECT().Delete(gomock.Any(), tt.args.key).Return(tt.repositoryErr).Times(tt.serviceCallTimes) //nolint

			ruleService := &service.RuleService{RuleRepository: ruleRepositoryMock}

			err := ruleService.Delete(tt.args.ctx, tt.args.key) //nolint
			if (err != nil) != (tt.wantedErr != nil) {          //nolint
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantedErr != nil) //nolint
				return
			}

			if tt.wantedErr != nil { //nolint
				if !reflect.DeepEqual(tt.wantedErr, err) { //nolint
					t.Fatalf("Error is not the expected. Expected: %v - Actual: %v", tt.wantedErr, err) //nolint
				}
				return
			}
		})
	}
}