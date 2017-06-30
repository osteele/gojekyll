package liquid

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/osteele/gojekyll/config"
	lq "github.com/osteele/liquid"
	"github.com/osteele/liquid/expressions"
	"github.com/osteele/liquid/generics"
	"github.com/russross/blackfriday"
)

func AddJekyllFilters(engine lq.Engine, config config.Config) {
	// array filters
	engine.DefineFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	// TODO neither Liquid nor Jekyll docs this, but it appears to be present
	engine.DefineFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		out := []interface{}{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				out = append(out, value)
			}
		}
		return out
	})
	engine.DefineFilter("group_by", groupByFilter)
	// sort overrides the Liquid filter with one that takes parameters
	engine.DefineFilter("sort", sortFilter)
	engine.DefineFilter("where", whereFilter) // TODO test case
	engine.DefineFilter("where_exp", whereExpFilter)
	engine.DefineFilter("xml_escape", xml.Marshal)

	engine.DefineFilter("push", func(array []interface{}, item interface{}) interface{} {
		return append(array, generics.MustConvertItem(item, array))
	})
	engine.DefineFilter("unshift", func(array []interface{}, item interface{}) interface{} {
		return append([]interface{}{generics.MustConvertItem(item, array)}, array...)
	})

	// dates
	engine.DefineFilter("date_to_rfc822", func(date time.Time) string {
		return date.Format(time.RFC822)
		// Out: Mon, 07 Nov 2008 13:07:54 -0800
	})
	engine.DefineFilter("date_to_string", func(date time.Time) string {
		return date.Format("02 Jan 2006")
		// Out: 07 Nov 2008
	})
	engine.DefineFilter("date_to_long_string", func(date time.Time) string {
		return date.Format("02 January 2006")
		// Out: 07 November 2008
	})
	engine.DefineFilter("date_to_xmlschema", func(date time.Time) string {
		return date.Format("2006-01-02T15:04:05-07:00")
		// Out: 2008-11-07T13:07:54-08:00
	})

	// strings
	engine.DefineFilter("absolute_url", func(s string) string {
		return config.AbsoluteURL + config.BaseURL + s
	})
	engine.DefineFilter("relative_url", func(s string) string {
		return config.BaseURL + s
	})
	engine.DefineFilter("jsonify", json.Marshal)
	engine.DefineFilter("markdownify", blackfriday.MarkdownCommon)
	// engine.DefineFilter("normalize_whitespace", func(s string) string {
	// 	wsPattern := regexp.MustCompile(`(?s:[\s\n]+)`)
	// 	return wsPattern.ReplaceAllString(s, " ")
	// })
	engine.DefineFilter("to_integer", func(n int) int { return n })
	engine.DefineFilter("number_of_words", func(s string) int {
		wordPattern := regexp.MustCompile(`\w+`) // TODO what's the Jekyll spec for a word?
		m := wordPattern.FindAllStringIndex(s, -1)
		if m == nil {
			return 0
		}
		return len(m)
	})

	// string escapes
	// engine.DefineFilter("uri_escape", func(s string) string {
	// 	parts := strings.SplitN(s, "?", 2)
	// 	if len(parts) > 0 {
	// TODO PathEscape is the wrong function
	// 		parts[len(parts)-1] = url.PathEscape(parts[len(parts)-1])
	// 	}
	// 	return strings.Join(parts, "?")
	// })
	engine.DefineFilter("xml_escape", func(s string) string {
		// TODO can't handle maps
		// eval https://github.com/clbanning/mxj
		// adapt https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
		buf := new(bytes.Buffer)
		if err := xml.EscapeText(buf, []byte(s)); err != nil {
			panic(err)
		}
		return buf.String()
	})
}

func arrayToSentenceStringFilter(array []string, conjunction interface{}) string {
	conj, ok := conjunction.(string)
	if !ok {
		conj = "and "
	}
	rt := reflect.ValueOf(array)
	ar := make([]string, rt.Len())
	for i, v := range array {
		ar[i] = v
		if i == rt.Len()-1 {
			ar[i] = conj + v
		}
	}
	return strings.Join(ar, ", ")
}

func groupByFilter(array []map[string]interface{}, property string) []map[string]interface{} {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil
	}
	groups := map[interface{}][]interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(property))
			// fmt.Println("invalid", item)
			if attr.IsValid() {
				key := attr.Interface()
				group, found := groups[key]
				// fmt.Println("found", attr)
				if found {
					group = append(group, groups[key])
				} else {
					group = []interface{}{item}
				}
				groups[key] = group
			}
		}
	}
	out := []map[string]interface{}{}
	for k, v := range groups {
		out = append(out, map[string]interface{}{"name": k, "items": v})
	}
	return out
}

func sortFilter(array []interface{}, key interface{}, nilFirst interface{}) []interface{} {
	nf, ok := nilFirst.(bool)
	if !ok {
		nf = true
	}
	out := make([]interface{}, len(array))
	copy(out, array)
	if key == nil {
		generics.Sort(out)
	} else {
		generics.SortByProperty(out, key.(string), nf)
	}
	return out
}

func whereExpFilter(array []interface{}, name string, expr expressions.Closure) ([]interface{}, error) {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, nil
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

func whereFilter(array []map[string]interface{}, key string, value interface{}) []interface{} {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil
	}
	out := []interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(key))
			if attr.IsValid() && fmt.Sprint(attr) == value {
				out = append(out, item.Interface())
			}
		}
	}
	return out
}
