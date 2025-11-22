package filters

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/osteele/gojekyll/config"
	"github.com/osteele/liquid"
	"github.com/stretchr/testify/require"
)

var filterTests = []struct{ in, expected string }{
	// dates
	// FIXME date_to_xmlschema should use local timezone?
	{`{{ time | date_to_xmlschema }}`, "2008-11-07T13:07:54+00:00"},
	{`{{ time | date_to_rfc822 }}`, "07 Nov 08 13:07 UTC"},
	{`{{ time | date_to_string }}`, "07 Nov 2008"},
	{`{{ time | date_to_long_string }}`, "07 November 2008"},

	// arrays
	{`{{ array | array_to_sentence_string }}`, "first, second, and third"},

	{`{{ site.members | group_by: "graduation_year" | map: "name" | sort | join }}`, "2013 2014 2015"},
	{`{{ site.members | group_by_exp: "item", "item.graduation_year" | size }}`, "4"},

	// TODO what is the default for nil first?
	{`{{ animals | sort | join: ", " }}`, "Sally Snake, giraffe, octopus, zebra"},
	{`{{ site.pages | sort: "weight" | map: "name" | join }}`, "b a d c"},
	{`{{ site.pages | sort: "weight", true | map: "name" | join }}`, "b a d c"},
	{`{{ site.pages | sort: "weight", false | map: "name" | join }}`, "a d c b"},

	{`{{ site.members | where: "graduation_year", "2014" | map: "name" }}`, "Alan"},
	{`{{ site.members | where_exp: "item", "item.graduation_year == 2014" | map: "name" }}`, "Alan"},
	{`{{ site.members | where_exp: "item", "item.graduation_year < 2014" | map: "name" }}`, "Alonzo"},
	{`{{ site.members | where_exp: "item", "item.name contains 'Al'" | map: "name" | join }}`, "Alonzo Alan"},

	{`{{ page.tags | push: 'Spokane' | join }}`, "Seattle Tacoma Spokane"},
	{`{{ page.tags | pop }}`, "Seattle"},
	{`{{ page.tags | shift }}`, "Tacoma"},
	{`{{ page.tags | unshift: "Olympia" | join }}`, "Olympia Seattle Tacoma"},

	// strings
	{`{{ "/assets/style.css" | relative_url }}`, "/my-baseurl/assets/style.css"},
	{`{{ "/assets/style.css" | absolute_url }}`, "http://example.com/my-baseurl/assets/style.css"},
	{`{{ "Markdown with _emphasis_ and *bold*." | markdownify }}`, "<p>Markdown with <em>emphasis</em> and <em>bold</em>.</p>"},
	{`{{ obj | jsonify }}`, `{"a":[1,2,3,4]}`},
	{`{{ site.pages | map: "name" | join }}`, "a b c d"},
	{`{{ site.pages | filter: "weight" | map: "name" | join }}`, "a c d"},
	{"{{ ws | normalize_whitespace }}", "a b c"},
	{`{{ "123" | to_integer | type }}`, "int"},
	{`{{ false | to_integer }}`, "0"},
	{`{{ true | to_integer }}`, "1"},
	{`{{ "here are some words" | number_of_words}}`, "4"},

	{`{{ "The _config.yml file" | slugify }}`, "the-config-yml-file"},
	{`{{ "The _config.yml file" | slugify: 'none' }}`, "the _config.yml file"},
	{`{{ "The _config.yml file" | slugify: 'raw' }}`, "the-_config.yml-file"},
	{`{{ "The _config.yml file" | slugify: 'default' }}`, "the-config-yml-file"},
	{`{{ "The _config.yml file" | slugify: 'pretty' }}`, "the-_config.yml-file"},

	// {`{{ "nav\n\tmargin: 0" | sassify }}`, "nav {\n  margin: 0; }"},
	{`{{ "nav {margin: 0}" | scssify }}`, "nav {\n  margin: 0;\n}"},

	{`{{ "smartify single 'quotes' here" | smartify }}`, "smartify single ‘quotes’ here"},
	{`{{ 'smartify double "quotes" here' | smartify }}`, "smartify double “quotes” here"},
	{"{{ \"smartify ``backticks''\" | smartify }}", "smartify “backticks”"},
	{`{{ "smartify it's they're" | smartify }}`, "smartify it’s they’re"},
	{`{{ "smartify ... (c) (r) (tm) -- ---" | smartify }}`, "smartify … © ® ™ – —"},

	{`{{ "foo, bar; baz?" | cgi_escape }}`, "foo%2C+bar%3B+baz%3F"},
	{`{{ "1 < 2 & 3" | xml_escape }}`, "1 &lt; 2 &amp; 3"},

	// Jekyll produces the first. I believe the second is acceptable.
	// {`{{ "http://foo.com/?q=foo, \bar?" | uri_escape }}`, "http://foo.com/?q=foo,%20%5Cbar?"},
	{`{{ "http://foo.com/?q=foo, \bar?" | uri_escape }}`, "http://foo.com/?q=foo%2C+%5Cbar%3F"},
}

var filterTestBindings = liquid.Bindings{
	"animals": []string{"zebra", "octopus", "giraffe", "Sally Snake"},
	"array":   []string{"first", "second", "third"},
	"obj": map[string]interface{}{
		"a": []int{1, 2, 3, 4},
	},
	"page": map[string]interface{}{
		"tags": []string{"Seattle", "Tacoma"},
	},
	"site": map[string]interface{}{
		"members": []map[string]interface{}{
			{"name": "Alonzo", "graduation_year": 2013},
			{"name": "Alan", "graduation_year": 2014},
			{"name": "Moses", "graduation_year": 2015},
			{"name": "Haskell"},
		},
		"pages": []map[string]interface{}{
			{"name": "a", "weight": 10},
			{"name": "b"},
			{"name": "c", "weight": 50},
			{"name": "d", "weight": 30},
		},
	},
	"time": timeMustParse("2008-11-07T13:07:54Z"),
	"ws":   "a  b\n\t c",
}

func TestFilters(t *testing.T) {
	for i, test := range filterTests {
		t.Run(fmt.Sprintf("%02d", i+1), func(t *testing.T) {
			requireTemplateRender(t, test.in, filterTestBindings, test.expected)
		})
	}
}

func TestSampleFilter(t *testing.T) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	AddJekyllFilters(engine, &cfg)

	// Test that sample returns one of the array elements
	data, err := engine.ParseAndRender([]byte(`{{ array | sample }}`), filterTestBindings)
	require.NoError(t, err)
	result := strings.TrimSpace(string(data))

	validResults := []string{"first", "second", "third"}
	require.Contains(t, validResults, result, "sample should return one of the array elements")
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

func TestWhereExpFilter_objects(t *testing.T) {
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

func requireTemplateRender(t *testing.T, tmpl string, bindings liquid.Bindings, expected string) {
	engine := liquid.NewEngine()
	cfg := config.Default()
	cfg.BaseURL = "/my-baseurl"
	cfg.AbsoluteURL = "http://example.com"
	AddJekyllFilters(engine, &cfg)
	data, err := engine.ParseAndRender([]byte(tmpl), bindings)
	require.NoErrorf(t, err, tmpl)
	require.Equalf(t, expected, strings.TrimSpace(string(data)), tmpl)
}

func timeMustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
