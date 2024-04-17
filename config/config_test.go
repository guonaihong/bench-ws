package config

import (
	"reflect"
	"testing"
)

func TestGenerateAddrs(t *testing.T) {
	type args struct {
		WSAddr string
		Name   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "TestGenerateAddrs",
			args: args{
				WSAddr: "ws://127.0.0.1:23001-23002/autobah",
				Name:   "",
			},
			want: []string{
				"ws://127.0.0.1:23001/autobah",
				"ws://127.0.0.1:23002/autobah",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateAddrs(tt.args.WSAddr, tt.args.Name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateAddrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
