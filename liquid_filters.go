package main

import (
	"reflect"

	"github.com/acstech/liquid/core"
)

func init() {
	core.RegisterFilter("where_exp", WhereExpFactory)
}

// WhereExpFactory implements the Jekyll `where_exp` filter
func WhereExpFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 2 {
		panic("were_exp requires two parameters")
	}
	// itemName := parameters[0]
	return (&WhereExpFilter{parameters[0], parameters[1]}).Run
}

// WhereExpFilter implements the Jekyll `where_exp` filter
type WhereExpFilter struct {
	varName core.Value
	expr    core.Value
}

// Run implements the Jekyll `where_exp` filter
func (f *WhereExpFilter) Run(input interface{}, data map[string]interface{}) interface{} {
	rt := reflect.ValueOf(input)
	switch rt.Kind() {
	case reflect.Slice:
	case reflect.Array:
	default:
		return input
	}

	varName := f.varName.Resolve(data).(string)
	expr := f.expr.Resolve(data).(string)
	p := core.NewParser([]byte(expr + "%}"))
	condition, err := p.ReadConditionGroup()
	// TODO assert we're at the end of the string
	if err != nil {
		panic(err)
	}

	result := []interface{}{}
	d := make(map[string]interface{})
	for k, v := range data {
		d[k] = v
	}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		d[varName] = item
		if condition.IsTrue(d) {
			result = append(result, item)
		}
	}
	return result
}
