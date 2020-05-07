package stringutils_test

import (
	"testing"

	stringutils "github.com/nicopozo/mockserver/internal/utils/string"
)

func TestArrayContains(t *testing.T) {
	type args struct {
		array []string
		str   string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Array should contain string",
			args: args{
				array: []string{"hola", "si", "tal vez", "otro string"},
				str:   "hola",
			},
			want: true,
		},
		{
			name: "Array should not contain string",
			args: args{
				array: []string{"hola", "si", "tal vez", "otro string"},
				str:   "chau",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutils.ArraysContains(tt.args.array, tt.args.str); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHash(t *testing.T) {
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
			got := stringutils.Hash(tt.args.s)
			if got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
