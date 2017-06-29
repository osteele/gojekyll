package liquid

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var filterTests = []struct{ in, expected string }{
	// dates
	// TODO date_to_xmlschema use local timezone
	{`{{time | date_to_xmlschema}}`, "2008-11-07T13:07:54+00:00"},
	{`{{time | date_to_rfc822 }}`, "07 Nov 08 13:07 UTC"},
	{`{{time | date_to_string }}`, "07 Nov 2008"},
	{`{{time | date_to_long_string }}`, "07 November 2008"},

	// arrays
	// TODO sort where group_by group_by_exp sample push pop shift
	// site.members | where:"graduation_year","2014"
	{`{{ar | array_to_sentence_string }}`, "first, second, and third"},

	// strings
	// TODO xml_escape cgi_escape uri_escape number_of_words
	// TODO scssify smartify slugify normalize_whitespace to_integer
	{`{{ "/assets/style.css" | relative_url }}`, "/my-baseurl/assets/style.css"},
	{`{{ "/assets/style.css" | absolute_url }}`, "http://example.com/my-baseurl/assets/style.css"},
	{`{{"Markdown with _emphasis_ and *bold*." | markdownify}}`, "<p>Markdown with <em>emphasis</em> and <em>bold</em>.</p>"},
	{`{{obj | jsonify }}`, `{"a":[1,2,3,4]}`},
	{`{{pages | map: "name" | join}}`, "a, b, c, d"},
	{`{{pages | filter: "weight" | map: "name" | join}}`, "a, c, d"},
}

var filterTestScope = map[string]interface{}{
	"ar": []string{"first", "second", "third"},
	"obj": map[string]interface{}{
		"a": []int{1, 2, 3, 4},
	},
	"pages": []map[string]interface{}{
		{"name": "a", "weight": 10},
		{"name": "b"},
		{"name": "c", "weight": 50},
		{"name": "d", "weight": 30},
	},
	"time": timeMustParse("2008-11-07T13:07:54Z"),
}

func timeMustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestFilters(t *testing.T) {
	for i, test := range filterTests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			requireTemplateRender(t, test.in, filterTestScope, test.expected)
		})
	}
}

func requireTemplateRender(t *testing.T, tmpl string, scope map[string]interface{}, expected string) {
	engine := NewEngine()
	engine.BaseURL = "/my-baseurl"
	engine.AbsoluteURL = "http://example.com"
	data, err := engine.ParseAndRender([]byte(tmpl), scope)
	require.NoErrorf(t, err, tmpl)
	require.Equalf(t, expected, strings.TrimSpace(string(data)), tmpl)
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
