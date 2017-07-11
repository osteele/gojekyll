package utils

import (
	"fmt"
	"reflect"

	"github.com/osteele/liquid"
)

// FollowDots applied to a property list ["a", "b", "c"] is equivalent to
// the Liquid data expression "data.a.b.c", except without special treatment
// of "first", "last", and "size".
func FollowDots(data interface{}, props []string) (interface{}, error) {
	for _, name := range props {
		data = liquid.FromDrop(data)
		if reflect.TypeOf(data).Kind() == reflect.Map {
			item := reflect.ValueOf(data).MapIndex(reflect.ValueOf(name))
			if item.IsValid() && !item.IsNil() && item.CanInterface() {
				data = item.Interface()
				continue
			}
		}
		return nil, fmt.Errorf("no such property: %q", name)
	}
	return data, nil
}
