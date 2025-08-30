package plugins

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
	"github.com/osteele/liquid/tags"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type jekyllSEOTagPlugin struct {
	plugin
	site Site
	tpl  *liquid.Template
}

func init() {
	register("jekyll-seo-tag", &jekyllSEOTagPlugin{})
}

func (p *jekyllSEOTagPlugin) AfterInitSite(s Site) error {
	p.site = s
	return nil
}

func (p *jekyllSEOTagPlugin) ConfigureTemplateEngine(e *liquid.Engine) error {
	e.RegisterTag("seo", p.seoTag)
	tpl, err := e.ParseTemplate([]byte(seoTagTemplateSource))
	if err != nil {
		panic(err)
	}
	p.tpl = tpl
	return nil
}

func (p *jekyllSEOTagPlugin) seoTag(ctx render.Context) (string, error) {
	buf := new(bytes.Buffer)
	e := seoTagTemplate.Execute(buf, seoTag{p.tpl, ctx})
	return buf.String(), e
}

type seoTag struct {
	tpl *liquid.Template
	ctx render.Context
}

var seoSiteFields = []string{"url", "twitter", "facebook", "logo", "social", "google_site_verification", "lang"}
var seoPageOrSiteFields = []string{"author", "description", "image", "author", "lang"}

func (p seoTag) TagBody() (string, error) {
	var (
		ctx          = p.ctx
		site         = liquid.FromDrop(ctx.Get("site")).(tags.IterationKeyedMap)
		page         = liquid.FromDrop(ctx.Get("page")).(tags.IterationKeyedMap)
		pageTitle    = page["title"]
		siteTitle    = site["title"]
		canonicalURL string
	)

	// Check if page has a custom canonical_url in front matter
	canonicalURL = fmt.Sprintf("%s%s", site["url"], page["url"])
	if pageCanonical := page["canonical_url"]; pageCanonical != nil {
		if s, ok := pageCanonical.(string); ok && s != "" {
			canonicalURL = s
		}
	}
	if siteTitle == nil {
		siteTitle = site["name"]
	}
	seoTag := map[string]interface{}{
		"title?": true,
		"title":  siteTitle,
		// the following are not doc'ed, but evident from inspection:
		// FIXME canonical w|w/out site.url and site.prefix
		"canonical_url": canonicalURL,
		"page_lang":     "en_US",
		"page_title":    pageTitle,
		"site_title":    siteTitle,
	}
	copyFields(seoTag, site, append(seoSiteFields, seoPageOrSiteFields...))
	copyFields(seoTag, page, seoPageOrSiteFields)
	// seoTag["description"] = page["excerpt"] FIXME appears to behave something like following, but excerpt should not be processed
	if pageTitle != nil && siteTitle != nil && pageTitle != siteTitle {
		seoTag["title"] = fmt.Sprintf("%s | %s", pageTitle, siteTitle)
	}
	if author, ok := seoTag["author"].(string); ok {
		if data, _ := utils.FollowDots(site, []string{"data", "authors", author}); data != nil {
			seoTag["author"] = data
		}
	}
	seoTag["json_ld"] = makeJSONLD(seoTag)
	bindings := map[string]interface{}{
		"page":    page,
		"site":    site,
		"seo_tag": seoTag,
	}
	b, err := p.tpl.Render(bindings)
	if err != nil {
		return "", err
	}
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	min := bytes.NewBuffer(make([]byte, 0, len(b)))
	if err := m.Minify("text/html", min, bytes.NewBuffer(b)); err != nil {
		return "", err
	}
	return min.String(), nil
}

func copyFields(to, from map[string]interface{}, fields []string) {
	for _, name := range fields {
		if value := from[name]; value != nil {
			to[name] = value
		}
	}
}

func makeJSONLD(seoTag map[string]interface{}) interface{} {
	jsonLD := map[string]interface{}{
		// TODO dateModified, datePublished, publisher, mainEntityOfPage
		"@context":    "http://schema.org",
		"@type":       "WebPage", // FIXME this is BlogPosting for a blog page
		"headline":    seoTag["page_title"],
		"description": seoTag["description"],
		"url":         seoTag["canonical_url"],
	}
	if author := seoTag["author"]; author != nil {
		if m, ok := author.(map[string]interface{}); ok {
			author = m["name"]
		}
		jsonLD["author"] = map[string]interface{}{
			"@type": "Person",
			"name":  author,
		}
	}
	return jsonLD
}

// This is a separate template so it isn't minimized away.
var seoTagTemplate = template.Must(template.New("SEO tag").Parse(
	`<!-- Begin Jekyll SEO tag -->
{{.TagBody}}
<!-- End Jekyll SEO tag -->`))

// Taken verbatim from https://github.com/jekyll/jekyll-seo-tag/
// Adapted from github.com/jekyll/jekyll-seo-tag. Used according to the MIT License.
const seoTagTemplateSource = `{% if seo_tag.title? %}
  <title>{{ seo_tag.title }}</title>
{% endif %}

{% if seo_tag.page_title %}
  <meta property="og:title" content="{{ seo_tag.page_title }}" />
{% endif %}

{% if seo_tag.author.name %}
  <meta name="author" content="{{ seo_tag.author.name }}" />
{% endif %}

<meta property="og:locale" content="{{ seo_tag.page_lang | replace:'-','_' }}" />

{% if seo_tag.description %}
  <meta name="description" content="{{ seo_tag.description }}" />
  <meta property="og:description" content="{{ seo_tag.description }}" />
{% endif %}

{% if site.url %}
  <link rel="canonical" href="{{ seo_tag.canonical_url }}" />
  <meta property="og:url" content="{{ seo_tag.canonical_url }}" />
{% endif %}

{% if seo_tag.site_title %}
  <meta property="og:site_name" content="{{ seo_tag.site_title }}" />
{% endif %}

{% if seo_tag.image %}
  <meta property="og:image" content="{{ seo_tag.image.path }}" />
  {% if seo_tag.image.height %}
    <meta property="og:image:height" content="{{ seo_tag.image.height }}" />
  {% endif %}
  {% if seo_tag.image.width %}
    <meta property="og:image:width" content="{{ seo_tag.image.width }}" />
  {% endif %}
{% endif %}

{% if page.date %}
  <meta property="og:type" content="article" />
  <meta property="article:published_time" content="{{ page.date | date_to_xmlschema }}" />
{% endif %}

{% if paginator.previous_page %}
  <link rel="prev" href="{{ paginator.previous_page_path | absolute_url }}">
{% endif %}
{% if paginator.next_page %}
  <link rel="next" href="{{ paginator.next_page_path | absolute_url }}">
{% endif %}

{% if site.twitter %}
  {% if seo_tag.image %}
    <meta name="twitter:card" content="summary_large_image" />
  {% else %}
    <meta name="twitter:card" content="summary" />
  {% endif %}

  <meta name="twitter:site" content="@{{ site.twitter.username | replace:"@","" }}" />

  {% if seo_tag.author.twitter %}
    <meta name="twitter:creator" content="@{{ seo_tag.author.twitter }}" />
  {% endif %}
{% endif %}

{% if site.facebook %}
  {% if site.facebook.admins %}
    <meta property="fb:admins" content="{{ site.facebook.admins }}" />
  {% endif %}

  {% if site.facebook.publisher %}
    <meta property="article:publisher" content="{{ site.facebook.publisher }}" />
  {% endif %}

  {% if site.facebook.app_id %}
    <meta property="fb:app_id" content="{{ site.facebook.app_id }}" />
  {% endif %}
{% endif %}

{% if site.webmaster_verifications %}
  {% if site.webmaster_verifications.google %}
    <meta name="google-site-verification" content="{{ site.webmaster_verifications.google }}">
  {% endif %}

  {% if site.webmaster_verifications.bing %}
    <meta name="msvalidate.01" content="{{ site.webmaster_verifications.bing }}">
  {% endif %}

  {% if site.webmaster_verifications.alexa %}
    <meta name="alexaVerifyID" content="{{ site.webmaster_verifications.alexa }}">
  {% endif %}

  {% if site.webmaster_verifications.yandex %}
    <meta name="yandex-verification" content="{{ site.webmaster_verifications.yandex }}">
  {% endif %}
{% elsif site.google_site_verification %}
  <meta name="google-site-verification" content="{{ site.google_site_verification }}" />
{% endif %}

<script type="application/ld+json">
  {{ seo_tag.json_ld | jsonify }}
</script>`
