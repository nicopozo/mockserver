package jsonutils_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	jsonutils "github.com/nicopozo/mockserver/internal/utils/json"
)

type dto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		body io.Reader
	}

	tests := []struct {
		name    string
		args    args
		wanted  *dto
		wantErr bool
	}{
		{
			name:    "Unmarshal successfully",
			args:    args{body: strings.NewReader("{\"id\" : 12,\"name\":\"the_name\"}")},
			wanted:  &dto{ID: 12, Name: "the_name"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &dto{}

			if err := jsonutils.Unmarshal(tt.args.body, result); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(result, tt.wanted) {
				t.Errorf("Unmarshal() = %v, want %v", result, tt.wanted)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	type args struct {
		model interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Marshal successfully",
			args: args{model: &dto{ID: 12, Name: "the_name"}},
			want: "{\"id\":12,\"name\":\"the_name\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonutils.Marshal(tt.args.model); got != tt.want {
				t.Errorf("Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
