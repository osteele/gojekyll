package liquid

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/acstech/liquid"
	"github.com/stretchr/testify/require"
)

var filterTests = []struct{ in, expected string }{
	{"{{time | date_to_rfc822 }}", "02 Jan 06 15:04 UTC"},
	{"{{obj | jsonify }}", `{"a":[1,2,3,4]}`},
	{"{{ar | array_to_sentence_string }}", "first, second, and third"},
}

var filterTestScope = map[string]interface{}{
	"ar": []string{"first", "second", "third"},
	"obj": map[string]interface{}{
		"a": []int{1, 2, 3, 4},
	},
	"time": timeMustParse("2006-01-02T15:04:05Z"),
}

func timeMustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		panic(err)
	}
	return t
}

func TestFilters(t *testing.T) {
	for _, test := range filterTests {
		requireTemplateRender(t, test.in, filterTestScope, test.expected)
	}
}

func requireTemplateRender(t *testing.T, tmpl string, data map[string]interface{}, expected string) {
	template, err := liquid.ParseString(tmpl, nil)
	require.NoError(t, err)
	writer := new(bytes.Buffer)
	template.Render(writer, data)
	require.Equal(t, expected, strings.TrimSpace(writer.String()))
}

// func TestXMLEscapeFilter(t *testing.T) {
// 	data := map[string]interface{}{
// 		"obj": map[string]interface{}{
// 			"a": []int{1, 2, 3, 4},
// 		},
// 	}
// 	requireTemplateRender(t, `{{obj | xml_escape }}`, data, `{"ak":[1,2,3,4]}`)
// }

func TestWhereExpFilter(t *testing.T) {
	var tmpl = `
	{% assign filtered = array | where_exp: "n", "n > 2" %}
	{% for item in filtered %}{{item}}{% endfor %}
	`
	data := map[string]interface{}{
		"array": []int{1, 2, 3, 4},
	}
	requireTemplateRender(t, tmpl, data, "34")
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
	requireTemplateRender(t, tmpl, data, "A")
}
