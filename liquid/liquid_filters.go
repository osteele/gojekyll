package liquid

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"time"

	"github.com/acstech/liquid"
	"github.com/acstech/liquid/core"
)

func init() {
	liquid.Tags["link"] = LinkFactory
	core.RegisterFilter("date_to_rfc822", DateToRFC822Factory)
	core.RegisterFilter("jsonify", JsonifyFactory)
	core.RegisterFilter("xml_escape", XMLEscapeFactory)
	core.RegisterFilter("where_exp", WhereExpFactory)
}

// DateToRFC822Factory implements the Jekyll `json` filter
func DateToRFC822Factory(parameters []core.Value) core.Filter {
	if len(parameters) != 0 {
		panic("The date_to_rfc822 filter doesn't accept parameters")
	}
	return func(input interface{}, data map[string]interface{}) interface{} {
		date := input.(time.Time) // TODO if a string, try parsing it
		return date.Format(time.RFC822)
	}
}

// JsonifyFactory implements the Jekyll `json` filter
func JsonifyFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 0 {
		panic("The jsonify filter doesn't accept parameters")
	}
	return func(input interface{}, data map[string]interface{}) interface{} {
		s, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}
		return s
	}
}

// XMLEscapeFactory implements the Jekyll `xml_escape` filter
func XMLEscapeFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 0 {
		panic("The xml_escape filter doesn't accept parameters")
	}
	return func(input interface{}, data map[string]interface{}) interface{} {
		s, err := xml.Marshal(input)
		// TODO can't handle maps
		// eval https://github.com/clbanning/mxj
		// adapt https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
		if err != nil {
			panic(err)
		}
		return s
	}
}

// WhereExpFactory implements the Jekyll `where_exp` filter
func WhereExpFactory(parameters []core.Value) core.Filter {
	if len(parameters) != 2 {
		panic("The were_exp filter requires two parameters")
	}
	return (&whereExpFilter{parameters[0], parameters[1]}).run
}

type whereExpFilter struct {
	varName core.Value
	expr    core.Value
}

func (f *whereExpFilter) run(input interface{}, data map[string]interface{}) interface{} {
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
