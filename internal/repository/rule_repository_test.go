package repository

import (
	"testing"
)

func Test_createExpression(t *testing.T) {
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

	for _, tt := range tests { //nolint
		t.Run(tt.name, func(t *testing.T) { //nolint
			if got := createExpression(tt.args.path); got != tt.want { //nolint
				t.Errorf("createExpression() = %v, want %v", got, tt.want) //nolint
			}
		})
	}
}
