package filters

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	sass "github.com/bep/godartsass/v2"
	"github.com/danog/blackfriday/v2"
	"github.com/osteele/gojekyll/config"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/evaluator"
	"github.com/osteele/liquid/expressions"
)

// AddJekyllFilters adds the Jekyll filters to the Liquid engine.
func AddJekyllFilters(e *liquid.Engine, c *config.Config) {
	// array filters
	e.RegisterFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	// TODO doc neither Liquid nor Jekyll docs this, but it appears to be present
	e.RegisterFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		var result []interface{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				result = append(result, value)
			}
		}
		return result
	})
	e.RegisterFilter("group_by", groupByFilter)
	e.RegisterFilter("group_by_exp", groupByExpFilter)
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
	e.RegisterFilter("pop", requireNonEmptyArray(func(array []interface{}) interface{} {
		return array[0]
	}))
	e.RegisterFilter("shift", requireNonEmptyArray(func(array []interface{}) interface{} {
		return array[len(array)-1]
	}))
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
		return utils.URLJoin(c.AbsoluteURL, c.BaseURL, s)
	})
	e.RegisterFilter("relative_url", func(s string) string {
		return c.BaseURL + s
	})
	e.RegisterFilter("jsonify", json.Marshal)
	e.RegisterFilter("markdownify", blackfriday.Run)
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
	e.RegisterFilter("cgi_escape", url.QueryEscape)
	e.RegisterFilter("sassify", unimplementedFilter("sassify"))
	e.RegisterFilter("scssify", scssifyFilter)
	e.RegisterFilter("smartify", smartifyFilter)
	e.RegisterFilter("uri_escape", func(s string) string {
		return regexp.MustCompile(`\?(.+?)=([^&]*)(?:\&(.+?)=([^&]*))*`).ReplaceAllStringFunc(s, func(m string) string {
			pair := strings.SplitN(m, "=", 2)
			return pair[0] + "=" + url.QueryEscape(pair[1])
		})
	})
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

// helpers

func requireNonEmptyArray(fn func([]interface{}) interface{}) func([]interface{}) interface{} {
	return func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return fn(array)
	}
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

// array filters

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

func groupByExpFilter(array []map[string]interface{}, name string, expr expressions.Closure) ([]map[string]interface{}, error) {
	rt := reflect.ValueOf(array)
	if !(rt.Kind() != reflect.Array || rt.Kind() == reflect.Slice) {
		return nil, nil
	}
	groups := map[interface{}][]interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		key, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if group, found := groups[key]; found {
			groups[key] = append(group, item)
		} else {
			groups[key] = []interface{}{item}
		}
	}
	var result []map[string]interface{}
	for k, v := range groups {
		result = append(result, map[string]interface{}{"name": k, "items": v})
	}
	return result, nil
}

func groupByFilter(array []map[string]interface{}, property string) []map[string]interface{} {
	rt := reflect.ValueOf(array)
	if !(rt.Kind() != reflect.Array || rt.Kind() == reflect.Slice) {
		return nil
	}
	groups := map[interface{}][]interface{}{}
	for i := 0; i < rt.Len(); i++ {
		irt := rt.Index(i)
		if irt.Kind() == reflect.Map && irt.Type().Key().Kind() == reflect.String {
			krt := irt.MapIndex(reflect.ValueOf(property))
			if krt.IsValid() && krt.CanInterface() {
				key := krt.Interface()
				if group, found := groups[key]; found {
					groups[key] = append(group, irt.Interface())
				} else {
					groups[key] = []interface{}{irt.Interface()}
				}
			}
		}
	}
	var result []map[string]interface{}
	for k, v := range groups {
		result = append(result, map[string]interface{}{"name": k, "items": v})
	}
	return result
}

func sortFilter(array []interface{}, key interface{}, nilFirst func(bool) bool) []interface{} {
	nf := nilFirst(true)
	result := make([]interface{}, len(array))
	copy(result, array)
	if key == nil {
		evaluator.Sort(result)
	} else {
		// TODO error if key is not a string
		evaluator.SortByProperty(result, key.(string), nf)
	}
	return result
}

func whereExpFilter(array []interface{}, name string, expr expressions.Closure) ([]interface{}, error) {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, nil
	}
	var result []interface{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		value, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if value != nil && value != false {
			result = append(result, item)
		}
	}
	return result, nil
}

func whereFilter(array []map[string]interface{}, key string, value interface{}) []interface{} {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil
	}
	var result []interface{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(key))
			if attr.IsValid() && fmt.Sprint(attr) == value {
				result = append(result, item.Interface())
			}
		}
	}
	return result
}

// string filters
var comp, compErr = sass.Start(sass.Options{})

func scssifyFilter(s string) (string, error) {
	if compErr != nil {
		return "", compErr
	}
	res, err := comp.Execute(sass.Args{
		Source: s,
	})
	return res.CSS, err
}
