package agent

import (
	"reflect"
	"runtime"
	"testing"
)

func TestGetMemStatByName(t *testing.T) {
	type args struct {
		a      *runtime.MemStats
		metric string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 reflect.Kind
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetMemStatByName(tt.args.a, tt.args.metric)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMemStatByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetMemStatByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
