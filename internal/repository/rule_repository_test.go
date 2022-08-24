package repository_test

import (
	"testing"

	"github.com/nicopozo/mockserver/internal/repository"
	"github.com/stretchr/testify/assert"
)

func Test_CreateExpression(t *testing.T) {
	type args struct {
		path string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with one param in the middle",
			args: args{
				path: "/first/{param}/second",
			},
			want: "^/first/[^/]+?/second$",
		},
		{
			name: "with one param in the end",
			args: args{
				path: "/first/second/{param}",
			},
			want: "^/first/second/[^/]+$",
		},
		{
			name: "with one param in the middle and one int the end",
			args: args{
				path: "/first/{param}/second/{param2}",
			},
			want: "^/first/[^/]+?/second/[^/]+$",
		},
	}

	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got := repository.CreateExpression(tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}
