package metric

import (
	"reflect"
	"testing"
)

func TestMakeMetricStruct(t *testing.T) {
	type args struct {
		name  string
		typeM MetricType
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Metric
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeMetricStruct(tt.args.name, tt.args.typeM, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeMetricStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeMetricStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}
