package log_test

import (
	"testing"

	"github.com/nicopozo/mockserver/internal/utils/log"
)

//nolint:nosnakecase,funlen
func Test_log_getMessage(t *testing.T) {
	type fields struct {
		trackingID string
	}

	type args struct {
		message string
		args    []interface{}
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Without args",
			fields: fields{trackingID: "123"},
			args: args{
				message: "the log message",
				args:    nil,
			},
			want: "the log message",
		},
		{
			name:   "with one string argument",
			fields: fields{trackingID: "123"},
			args: args{
				message: "with one string %v message",
				args:    []interface{}{"argument"},
			},
			want: "with one string argument message",
		},
		{
			name:   "with one int argument",
			fields: fields{trackingID: "123"},
			args: args{
				message: "with %v int argument message",
				args:    []interface{}{1},
			},
			want: "with 1 int argument message",
		},
		{
			name:   "with one string argument and one int",
			fields: fields{trackingID: "123"},
			args: args{
				message: "with one string %s and %d int value message",
				args:    []interface{}{"argument", 1},
			},
			want: "with one string argument and 1 int value message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theLogger := log.NewLogger(tt.fields.trackingID)
			if got := theLogger.GetMessage(tt.args.message, tt.args.args...); got != tt.want {
				t.Errorf("log.GetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
