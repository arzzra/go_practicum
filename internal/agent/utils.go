package agent

import (
	"reflect"
	"runtime"
)

func GetMemStatByName(a *runtime.MemStats, metric string) (interface{}, reflect.Kind) {
	r := reflect.ValueOf(a)
	f := reflect.Indirect(r).FieldByName(metric)
	if f.IsValid() {
		typevalue := f.Type().Kind()
		switch typevalue {
		default:
			return nil, 0
		case reflect.Float64:
			return f.Float(), typevalue
		case reflect.Uint64:
			return f.Uint(), typevalue
		}
	}
	return nil, 0
}
