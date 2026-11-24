package stringutils_test

import (
	"testing"

	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
)

func TestHash(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Hash should generate correctly for string",
			args: args{s: "hello"},
			want: uint32(1335831723),
		},
		{
			name: "Hash should generate correctly for long string",
			args: args{s: "Lorem ipsum dolor sit amet, consectetur adipiscing elit"},
			want: uint32(2580452767),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := stringutils.Hash(tt.args.s)
			if got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
