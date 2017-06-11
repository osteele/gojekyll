package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/acstech/liquid"
)

func assertExpected(t *testing.T, actual, expected string) {
	if actual != expected {
		t.Fatalf(fmt.Sprintf("expected %s; got %s", expected, actual))
	}
}

func testTemplateRender(t *testing.T, tmpl string, data map[string]interface{}, expected string) {
	template, err := liquid.ParseString(tmpl, nil)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	writer := new(bytes.Buffer)
	template.Render(writer, data)
	assertExpected(t, expected, strings.TrimSpace(writer.String()))
}

func TestWhereExp(t *testing.T) {
	var tmpl = `
	{% assign filtered = array | where_exp: "n", "n > 2" %}
	{% for item in filtered %}{{item}}{% endfor %}
	`

	data := map[string]interface{}{
		"array": []int{1, 2, 3, 4},
	}

	testTemplateRender(t, tmpl, data, "34")
}

func TestWhereExpObjects(t *testing.T) {
	var tmpl = `
	{% assign filtered = array | where_exp: "item", "item.flag == true" %}
	{% for item in filtered %}{{item.name}}{% endfor %}
	`

	data := map[string]interface{}{
		"array": []map[string]interface{}{
			map[string]interface{}{
				"name": "A",
				"flag": true,
			},
			map[string]interface{}{
				"name": "B",
				"flag": false,
			},
		}}

	testTemplateRender(t, tmpl, data, "A")
}
