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

// nolint:funlen,nosnakecase,paralleltest
func Test_ruleFileRepository_Get(t *testing.T) {
	type args struct {
		ctx context.Context //nolint:containedctx
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

	for _, tt := range tests { //nolint:paralleltest,varnamelen
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

func getMocksFile() string {
	dest, err := os.OpenFile("mocks.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600) //nolint:nosnakecase
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
