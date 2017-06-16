package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/acstech/liquid"
	"github.com/stretchr/testify/assert"
)

func assertTemplateRender(t *testing.T, tmpl string, data VariableMap, expected string) {
	template, err := liquid.ParseString(tmpl, nil)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	writer := new(bytes.Buffer)
	template.Render(writer, data)
	assert.Equal(t, expected, strings.TrimSpace(writer.String()))
}

func TestDateToRFC822Filter(t *testing.T) {
	t0, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
	data := VariableMap{"time": t0}
	assertTemplateRender(t, `{{time | date_to_rfc822 }}`, data, "02 Jan 06 15:04 UTC")
}

func TestJsonifyFilter(t *testing.T) {
	data := VariableMap{
		"obj": map[string]interface{}{
			"a": []int{1, 2, 3, 4},
		},
	}
	assertTemplateRender(t, `{{obj | jsonify }}`, data, `{"a":[1,2,3,4]}`)
}

// func TestXMLEscapeFilter(t *testing.T) {
// 	data := VariableMap{
// 		"obj": map[string]interface{}{
// 			"a": []int{1, 2, 3, 4},
// 		},
// 	}
// 	assertTemplateRender(t, `{{obj | xml_escape }}`, data, `{"ak":[1,2,3,4]}`)
// }

func TestWhereExpFilter(t *testing.T) {
	var tmpl = `
	{% assign filtered = array | where_exp: "n", "n > 2" %}
	{% for item in filtered %}{{item}}{% endfor %}
	`
	data := VariableMap{
		"array": []int{1, 2, 3, 4},
	}
	assertTemplateRender(t, tmpl, data, "34")
}

func TestWhereExpFilterObjects(t *testing.T) {
	var tmpl = `
	{% assign filtered = array | where_exp: "item", "item.flag == true" %}
	{% for item in filtered %}{{item.name}}{% endfor %}
	`
	data := map[string]interface{}{
		"array": []map[string]interface{}{
			{
				"name": "A",
				"flag": true,
			},
			{
				"name": "B",
				"flag": false,
			},
		}}
	assertTemplateRender(t, tmpl, data, "A")
}
