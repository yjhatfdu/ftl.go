package functions

import (
	"fmt"
	"reflect"
	"time"
)

var FuncMap = map[string]func(args ...interface{}) (interface{}, error){
	"eq": eq,
}

func buildBinaryTypeError(op string, v1, v2 interface{}) error {
	return fmt.Errorf("无法对 %s 和 %s 执行 %s 操作", reflect.TypeOf(v1).Name(), reflect.TypeOf(v2).Name(), op)
}

func eq(args ...interface{}) (interface{}, error) {
	l := args[0]
	r := args[1]
	if l == nil || r == nil {
		return nil, nil
	}
	switch lv := l.(type) {
	case uint64, uint32, uint16, uint8, int32, int16, int8, int, uint:
		return reflect.ValueOf(l).Int()==reflect.ValueOf(r).Int(),nil
	case float64:
		switch rv := r.(type) {
		case float64:
			return lv == rv, nil
		}
	case string:
		switch rv := r.(type) {
		case string:
			return lv == rv, nil
		}
	case bool:
		switch rv := r.(type) {
		case bool:
			return lv == rv, nil
		}
	case time.Time:
		switch rv := r.(type) {
		case time.Time:
			return lv.Equal(rv), nil
		}
	}
	return nil, buildBinaryTypeError("=", l, r)
}
