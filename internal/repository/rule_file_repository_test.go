package repository_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	ruleserrors "github.com/nicopozo/mockserver/internal/errors"
	"github.com/nicopozo/mockserver/internal/model"
	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/stretchr/testify/assert"
)

func Test_ruleFileRepository_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}

	tests := []struct {
		name      string
		args      args
		want      *model.Rule
		wantedErr error
	}{
		{
			name: "Should get rule successfully",
			args: args{
				ctx: context.Background(),
				key: "a1",
			},
			want: &model.Rule{
				Key:      "a1",
				Group:    "TestApp",
				Name:     "TestMock",
				Path:     "/test",
				Strategy: "normal",
				Method:   "DELETE",
				Status:   "enabled",
				Responses: []model.Response{
					{
						Body:        "{\"field\":\"value\"}",
						ContentType: "application/json",
						HTTPStatus:  http.StatusOK,
						Delay:       0,
						Scene:       "",
					},
				},
				Variables: []*model.Variable{},
			},
			wantedErr: nil,
		},
		{
			name: "Should return NotFoundError when ID not found",
			args: args{
				ctx: context.Background(),
				key: "not-found-id",
			},
			want:      nil,
			wantedErr: ruleserrors.RuleNotFoundError{Message: "no rule found with key: not-found-id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := getMocksFile()
			defer func(fileName string) { _ = os.Remove(fileName) }(name)

			fileRepository, err := repository.NewRuleFileRepository(name)
			if assert.Nil(t, err) {
				got, err := fileRepository.Get(tt.args.ctx, tt.args.key)

				if tt.wantedErr != nil {
					assert.NotNil(t, err)
					assert.Equal(t, tt.wantedErr, err)

					return
				}

				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_ruleFileRepository_Search(t *testing.T) {
	type args struct {
		ctx    context.Context
		params map[string]interface{}
		paging model.Paging
	}

	tests := []struct {
		name    string
		args    args
		wantLen int
		total   int64
	}{
		{
			name: "Should search and return all rules",
			args: args{
				ctx:    context.Background(),
				params: map[string]interface{}{},
				paging: model.Paging{Limit: 10},
			},
			wantLen: 4,
			total:   4,
		},
		{
			name: "Should search with limit",
			args: args{
				ctx:    context.Background(),
				params: map[string]interface{}{},
				paging: model.Paging{Limit: 1},
			},
			wantLen: 1,
			total:   4,
		},
		{
			name: "Should search with LastID (keyset pagination)",
			args: args{
				ctx:    context.Background(),
				params: map[string]interface{}{},
				paging: model.Paging{Limit: 10, LastID: "c3"}, // c3 is index 1 after sorting newest first (d4, c3, b2, a1)
			},
			wantLen: 2, // b2, a1
			total:   4,
		},
		{
			name: "Should search with filters",
			args: args{
				ctx:    context.Background(),
				params: map[string]interface{}{"name": "TestMock"},
				paging: model.Paging{Limit: 10},
			},
			wantLen: 1,
			total:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := getMocksFile()
			defer func(fileName string) { _ = os.Remove(fileName) }(name)

			fileRepository, err := repository.NewRuleFileRepository(name)
			if assert.Nil(t, err) {
				got, err := fileRepository.Search(tt.args.ctx, tt.args.params, tt.args.paging)

				assert.Nil(t, err)
				assert.Equal(t, tt.wantLen, len(got.Results))
				assert.Equal(t, tt.total, got.Paging.Total)
			}
		})
	}
}

func getMocksFile() string {
	dest, err := os.OpenFile("mocks.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) { _ = f.Close() }(dest)

	originalFile, err := os.Open("../utils/test/mocks/json/mocks.json")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(dest, originalFile)
	if err != nil {
		panic(err)
	}

	return "mocks.json"
}
