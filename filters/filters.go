package filters

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/evaluator"
	"github.com/osteele/liquid/expression"
	"github.com/russross/blackfriday"
)

// AddJekyllFilters adds the Jekyll filters to the Liquid engine.
func AddJekyllFilters(e liquid.Engine, c config.Config) {
	// array filters
	e.RegisterFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	// TODO doc neither Liquid nor Jekyll docs this, but it appears to be present
	e.RegisterFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		out := []interface{}{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				out = append(out, value)
			}
		}
		return out
	})
	e.RegisterFilter("group_by", groupByFilter)
	e.RegisterFilter("group_by_exp", unimplementedFilter("group_by_exp"))
	e.RegisterFilter("sample", func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return array[rand.Intn(len(array))]
	})
	// sort overrides the Liquid filter with one that takes parameters
	e.RegisterFilter("sort", sortFilter)
	e.RegisterFilter("where", whereFilter) // TODO test case
	e.RegisterFilter("where_exp", whereExpFilter)
	e.RegisterFilter("xml_escape", xml.Marshal)
	e.RegisterFilter("push", func(array []interface{}, item interface{}) interface{} {
		return append(array, evaluator.MustConvertItem(item, array))
	})
	e.RegisterFilter("pop", func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return array[0]
	})
	e.RegisterFilter("shift", func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return array[len(array)-1]
	})
	e.RegisterFilter("unshift", func(array []interface{}, item interface{}) interface{} {
		return append([]interface{}{evaluator.MustConvertItem(item, array)}, array...)
	})

	// dates
	e.RegisterFilter("date_to_rfc822", func(date time.Time) string {
		return date.Format(time.RFC822)
		// Out: Mon, 07 Nov 2008 13:07:54 -0800
	})
	e.RegisterFilter("date_to_string", func(date time.Time) string {
		return date.Format("02 Jan 2006")
		// Out: 07 Nov 2008
	})
	e.RegisterFilter("date_to_long_string", func(date time.Time) string {
		return date.Format("02 January 2006")
		// Out: 07 November 2008
	})
	e.RegisterFilter("date_to_xmlschema", func(date time.Time) string {
		return date.Format("2006-01-02T15:04:05-07:00")
		// Out: 2008-11-07T13:07:54-08:00
	})

	// strings
	e.RegisterFilter("absolute_url", func(s string) string {
		return c.AbsoluteURL + c.BaseURL + s
	})
	e.RegisterFilter("relative_url", func(s string) string {
		return c.BaseURL + s
	})
	e.RegisterFilter("jsonify", json.Marshal)
	e.RegisterFilter("markdownify", blackfriday.MarkdownCommon)
	e.RegisterFilter("normalize_whitespace", func(s string) string {
		// s = strings.Replace(s, "n", "N", -1)
		wsPattern := regexp.MustCompile(`(?s:[\s\n]+)`)
		return wsPattern.ReplaceAllString(s, " ")
	})
	e.RegisterFilter("slugify", func(s, mode string) string {
		if mode == "" {
			mode = "default"
		}
		p := map[string]string{
			"raw":     `\s+`,
			"default": `[^[:alnum:]]+`,
			"pretty":  `[^[:alnum:]\._~!$&'()+,;=@]+`,
		}[mode]
		if p != "" {
			s = regexp.MustCompile(p).ReplaceAllString(s, "-")
		}
		return strings.ToLower(s)
	})
	e.RegisterFilter("to_integer", func(n int) int { return n })
	e.RegisterFilter("number_of_words", func(s string) int {
		wordPattern := regexp.MustCompile(`\w+`) // TODO what's the Jekyll spec for a word?
		m := wordPattern.FindAllStringIndex(s, -1)
		if m == nil {
			return 0
		}
		return len(m)
	})

	// string escapes
	// engine.RegisterFilter("uri_escape", func(s string) string {
	// 	parts := strings.SplitN(s, "?", 2)
	// 	if len(parts) > 0 {
	// TODO PathEscape is the wrong function
	// 		parts[len(parts)-1] = url.PathEscape(parts[len(parts)-1])
	// 	}
	// 	return strings.Join(parts, "?")
	// })
	e.RegisterFilter("cgi_escape", unimplementedFilter("cgi_escape"))
	e.RegisterFilter("uri_escape", unimplementedFilter("uri_escape"))
	e.RegisterFilter("scssify", unimplementedFilter("scssify"))
	e.RegisterFilter("smartify", unimplementedFilter("smartify"))
	e.RegisterFilter("xml_escape", func(s string) string {
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

func unimplementedFilter(name string) func(value interface{}) interface{} {
	warned := false
	return func(value interface{}) interface{} {
		if !warned {
			fmt.Println("warning: unimplemented filter:", name)
			warned = true
		}
		return value
	}
}

func arrayToSentenceStringFilter(array []string, conjunction func(string) string) string {
	conj := conjunction("and ")
	switch len(array) {
	case 1:
		return array[0]
	default:
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

func sortFilter(array []interface{}, key interface{}, nilFirst func(bool) bool) []interface{} {
	nf := nilFirst(true)
	out := make([]interface{}, len(array))
	copy(out, array)
	if key == nil {
		evaluator.Sort(out)
	} else {
		// TODO error if key is not a string
		evaluator.SortByProperty(out, key.(string), nf)
	}
	return out
}

func whereExpFilter(array []interface{}, name string, expr expression.Closure) ([]interface{}, error) {
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
