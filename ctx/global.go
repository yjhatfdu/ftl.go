package ctx

import (
	"fmt"
	"reflect"
	"strings"
)

type Global struct {
	Variables map[string]interface{}
	Root      interface{}
	Parent    *Global
}

func (g *Global) Get(key string) (interface{}, error) {
	//upperKey := strings.ToUpper(key[0:1]) + key[1:]
	if out, ok := g.Variables[key]; ok {
		return out, nil
	}
	//if out, ok := g.Variables[upperKey]; ok {
	//	return out, nil
	//}
	if g.Parent != nil {
		return g.Parent.Get(key)
	}
	return GetField(g.Root, key)
}

func (g *Global) Fork() *Global {
	return &Global{
		Variables: map[string]interface{}{},
		Parent:    g,
	}
}

func GetField(v interface{}, key string) (interface{}, error) {
	upperKey := strings.ToUpper(key[0:1]) + key[1:]
	if rootmap, ok := v.(map[string]interface{}); ok {
		if out, ok := rootmap[key]; ok {
			return out, nil
		}
		if out, ok := rootmap[upperKey]; ok {
			return out, nil
		}
	}
	r := reflect.ValueOf(v)
	if !r.IsValid() {
		return nil, nil
	}
	k := r.Kind()
	if k == reflect.Interface || k == reflect.Ptr {
		return GetField(r.Elem(), key)
	}
	if r.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot get field '%s' of kind %v", key, r.Kind())
	}
	out := r.FieldByName(key)
	if out.IsValid() {
		return out.Interface(), nil
	}
	out = r.FieldByName(upperKey)
	if out.IsValid() {
		return out.Interface(), nil
	}
	out = r.MethodByName(key)
	if out.IsValid() {
		return out.Interface(), nil
	}
	out = r.MethodByName(upperKey)
	if out.IsValid() {
		return out.Interface(), nil
	}
	return nil, nil
}
