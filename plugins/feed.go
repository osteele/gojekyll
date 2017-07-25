package plugins

import (
	"fmt"
	"html"
	"io"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

type jekyllFeedPlugin struct {
	plugin
	site Site
	tpl  *liquid.Template
}

func init() {
	register("jekyll-feed", &jekyllFeedPlugin{})
}

func (p *jekyllFeedPlugin) Initialize(s Site) error {
	p.site = s
	return nil
}

func (p *jekyllFeedPlugin) ConfigureTemplateEngine(e *liquid.Engine) error {
	e.RegisterTag("feed_meta", p.feedMetaTag)
	tpl, err := e.ParseTemplate([]byte(feedTemplateSource))
	if err != nil {
		panic(err)
	}
	p.tpl = tpl
	return nil
}

func (p *jekyllFeedPlugin) PostRead(s Site) error {
	path := "/feed.xml"
	if cfg, ok := s.Config().Variables["feed"].(map[string]interface{}); ok {
		if pp, ok := cfg["path"].(string); ok {
			path = "/" + pp
		}
	}
	d := feedDoc{s, p, path}
	s.AddDocument(&d, true)
	return nil
}

func (p *jekyllFeedPlugin) feedMetaTag(ctx render.Context) (string, error) {
	cfg := p.site.Config()
	name, _ := cfg.Variables["name"].(string)
	tag := fmt.Sprintf(`<link type="application/atom+xml" rel="alternate" href="%s/feed.xml" title="%s">`,
		html.EscapeString(cfg.AbsoluteURL), html.EscapeString(name))
	return tag, nil
}

type feedDoc struct {
	site   Site
	plugin *jekyllFeedPlugin
	path   string
}

func (d *feedDoc) Permalink() string  { return d.path }
func (d *feedDoc) SourcePath() string { return "" }
func (d *feedDoc) OutputExt() string  { return ".xml" }
func (d *feedDoc) Published() bool    { return true }
func (d *feedDoc) Static() bool       { return false } // FIXME means different things to different callers
func (d *feedDoc) Reload() error      { return nil }

func (d *feedDoc) Content() []byte {
	bindings := map[string]interface{}{"site": d.site}
	b, err := d.plugin.tpl.Render(bindings)
	if err != nil {
		panic(err)
	}
	return b
}

func (d *feedDoc) Write(w io.Writer) error {
	_, err := w.Write(d.Content())
	return err
}

// Taken verbatim from https://github.com/jekyll/jekyll-feed/
const feedTemplateSource = `<?xml version="1.0" encoding="utf-8"?>
{% if page.xsl %}
  <?xml-stylesheet type="text/xml" href="{{ '/feed.xslt.xml' | absolute_url }}"?>
{% endif %}
<feed xmlns="http://www.w3.org/2005/Atom" {% if site.lang %}xml:lang="{{ site.lang }}"{% endif %}>
  <generator uri="https://jekyllrb.com/" version="{{ jekyll.version }}">Jekyll</generator>
  <link href="{{ page.url | absolute_url }}" rel="self" type="application/atom+xml" />
  <link href="{{ '/' | absolute_url }}" rel="alternate" type="text/html" {% if site.lang %}hreflang="{{ site.lang }}" {% endif %}/>
  <updated>{{ site.time | date_to_xmlschema }}</updated>
  <id>{{ '/' | absolute_url | xml_escape }}</id>

  {% if site.title %}
    <title type="html">{{ site.title | smartify | xml_escape }}</title>
  {% elsif site.name %}
    <title type="html">{{ site.name | smartify | xml_escape }}</title>
  {% endif %}

  {% if site.description %}
    <subtitle>{{ site.description | xml_escape }}</subtitle>
  {% endif %}

  {% if site.author %}
    <author>
        <name>{{ site.author.name | default: site.author | xml_escape }}</name>
      {% if site.author.email %}
        <email>{{ site.author.email | xml_escape }}</email>
      {% endif %}
      {% if site.author.uri %}
        <uri>{{ site.author.uri | xml_escape }}</uri>
      {% endif %}
    </author>
  {% endif %}

  {% assign posts = site.posts | where_exp: "post", "post.draft != true" %}
  {% for post in posts limit: 10 %}
    <entry{% if post.lang %}{{" "}}xml:lang="{{ post.lang }}"{% endif %}>
      <title type="html">{{ post.title | smartify | strip_html | normalize_whitespace | xml_escape }}</title>
      <link href="{{ post.url | absolute_url }}" rel="alternate" type="text/html" title="{{ post.title | xml_escape }}" />
      <published>{{ post.date | date_to_xmlschema }}</published>
      <updated>{{ post.last_modified_at | default: post.date | date_to_xmlschema }}</updated>
      <id>{{ post.id | absolute_url | xml_escape }}</id>
      <content type="html" xml:base="{{ post.url | absolute_url | xml_escape }}">{{ post.content | strip | xml_escape }}</content>

      {% assign post_author = post.author | default: post.authors[0] | default: site.author %}
      {% assign post_author = site.data.authors[post_author] | default: post_author %}
      {% assign post_author_email = post_author.email | default: nil %}
      {% assign post_author_uri = post_author.uri | default: nil %}
      {% assign post_author_name = post_author.name | default: post_author %}

      <author>
          <name>{{ post_author_name | default: "" | xml_escape }}</name>
        {% if post_author_email %}
          <email>{{ post_author_email | xml_escape }}</email>
        {% endif %}
        {% if post_author_uri %}
          <uri>{{ post_author_uri | xml_escape }}</uri>
        {% endif %}
      </author>

      {% if post.category %}
        <category term="{{ post.category | xml_escape }}" />
      {% endif %}

      {% for tag in post.tags %}
        <category term="{{ tag | xml_escape }}" />
      {% endfor %}

      {% if post.excerpt and post.excerpt != empty %}
        <summary type="html">{{ post.excerpt | strip_html | normalize_whitespace | xml_escape }}</summary>
      {% endif %}

      {% assign post_image = post.image.path | default: post.image %}
      {% if post_image %}
        {% unless post_image contains "://" %}
          {% assign post_image = post_image | absolute_url | xml_escape  %}
        {% endunless %}
        <media:thumbnail xmlns:media="http://search.yahoo.com/mrss/" url="{{ post_image }}" />
      {% endif %}
    </entry>
  {% endfor %}
</feed>`
