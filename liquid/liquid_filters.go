package liquid

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"strings"
	"time"

	"github.com/osteele/liquid/expressions"
	"github.com/russross/blackfriday"
)

func (e *LocalWrapperEngine) addJekyllFilters() {
	// arrays
	e.engine.DefineFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	// TODO neither Liquid nor Jekyll docs this, but it appears to be present
	e.engine.DefineFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		out := []interface{}{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				out = append(out, value)
			}
		}
		return out
	})
	e.engine.DefineFilter("where_exp", whereExpFilter)
	e.engine.DefineFilter("xml_escape", xml.Marshal)

	// dates
	e.engine.DefineFilter("date_to_rfc822", func(date time.Time) string {
		return date.Format(time.RFC822)
		// Out: Mon, 07 Nov 2008 13:07:54 -0800
	})
	e.engine.DefineFilter("date_to_string", func(date time.Time) string {
		return date.Format("02 Jan 2006")
		// Out: 07 Nov 2008
	})
	e.engine.DefineFilter("date_to_long_string", func(date time.Time) string {
		return date.Format("02 January 2006")
		// Out: 07 November 2008
	})
	e.engine.DefineFilter("date_to_xmlschema", func(date time.Time) string {
		return date.Format("2006-01-02T15:04:05-07:00")
		// Out: 2008-11-07T13:07:54-08:00
	})

	// strings
	e.engine.DefineFilter("absolute_url", func(s string) string {
		return e.AbsoluteURL + e.BaseURL + s
	})
	e.engine.DefineFilter("relative_url", func(s string) string {
		return e.BaseURL + s
	})
	e.engine.DefineFilter("jsonify", func(value interface{}) []byte {
		s, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		return s
	})
	e.engine.DefineFilter("markdownify", blackfriday.MarkdownCommon)
}

func arrayToSentenceStringFilter(value []string, conjunction interface{}) string {
	conj, ok := conjunction.(string)
	if !ok {
		conj = "and "
	}
	rt := reflect.ValueOf(value)
	ar := make([]string, rt.Len())
	for i, v := range value {
		ar[i] = v
		if i == rt.Len()-1 {
			ar[i] = conj + v
		}
	}
	return strings.Join(ar, ", ")
}

// func xmlEscapeFilter(value interface{}) interface{} {
// 	data, err := xml.Marshal(value)
// 	// TODO can't handle maps
// 	// eval https://github.com/clbanning/mxj
// 	// adapt https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
// 	if err != nil {
// 		panic(err)
// 	}
// 	return data
// }

func whereExpFilter(in []interface{}, name string, expr expressions.Closure) ([]interface{}, error) {
	rt := reflect.ValueOf(in)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return in, nil
	}
	out := []interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		value, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if value != nil && value != false {
			out = append(out, item)
		}
	}
	return out, nil
}

func whereFilter(in []interface{}, key string, value interface{}) []interface{} {
	rt := reflect.ValueOf(in)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return in
	}
	out := []interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(key))
			if attr != reflect.Zero(attr.Type()) && (value == nil || attr.Interface() == value) {
				out = append(out, item.Interface())
			}
		}
	}
	return out
}
