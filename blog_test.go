package main

import (
	"encoding/json"
	"encoding/xml"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/ErikJermanis/eriksoftware.hr/blog"
)

func publishedFixture(slug, date string) blog.Post {
	p := blog.ExamplePost()
	p.Slug, p.Title, p.PublishedAt, p.Status = slug, "Članak "+slug, date, blog.Published
	return p
}

func fixtureHandler(t *testing.T, posts ...blog.Post) http.Handler {
	t.Helper()
	h, err := handlerWithPosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestExampleStaysOffPublicSite(t *testing.T) {
	h := testHandler(t)
	example := blog.ExamplePost()
	for _, path := range []string{example.Path(), example.Path() + "/"} {
		if r := request(t, h, path); r.Code != http.StatusNotFound {
			t.Errorf("example route %s = %d, want 404", path, r.Code)
		}
	}
	for _, path := range []string{"/blog", "/sitemap.xml", "/blog/feed.xml"} {
		if strings.Contains(request(t, h, path).Body.String(), example.Slug) {
			t.Errorf("example leaked into %s", path)
		}
	}
}

// New registrations are checked automatically; authors need not add a test for
// every article, and removing the example registration does not break tests.
func TestRegisteredPosts(t *testing.T) {
	h := testHandler(t)
	for _, p := range blog.Posts() {
		t.Run(p.Slug, func(t *testing.T) {
			r := request(t, h, p.Path())
			published, err := time.Parse(time.RFC3339, p.PublishedAt)
			if err != nil {
				t.Fatal(err)
			}
			if p.Status == blog.Draft || published.After(time.Now()) {
				if r.Code != 404 {
					t.Error("unpublished article is reachable")
				}
				if request(t, h, p.Path()+"/").Code != http.StatusNotFound {
					t.Error("unpublished article is reachable with a trailing slash")
				}
				for _, endpoint := range []string{"/blog", "/sitemap.xml", "/blog/feed.xml"} {
					if strings.Contains(request(t, h, endpoint).Body.String(), p.Slug) {
						t.Errorf("unpublished article leaked into %s", endpoint)
					}
				}
				return
			}
			page := r.Body.String()
			if r.Code != 200 || len(regexp.MustCompile(`<h1(?:\s|>)`).FindAllString(page, -1)) != 1 {
				t.Fatal("article must respond with 200 and exactly one H1")
			}
			if noIndex := strings.Contains(page, `content="noindex, follow"`); noIndex != (p.Status != blog.Published) {
				t.Error("incorrect publication/indexing state")
			}
			schemaNode(t, schemaGraph(t, page), "BlogPosting")
			assertLocalLinks(t, h, page)
		})
	}
}

func TestPublishedBlogListingAndRoutes(t *testing.T) {
	older := publishedFixture("stariji", "2024-01-02T09:00:00+01:00")
	newer := publishedFixture("noviji", "2024-03-02T09:00:00+01:00")
	draft := publishedFixture("nacrt", "2024-02-02T09:00:00+01:00")
	draft.Status = blog.Draft
	future := publishedFixture("buduci", "2099-01-02T09:00:00+01:00")
	h := fixtureHandler(t, older, draft, blog.ExamplePost(), future, newer)
	response := request(t, h, "/blog?utm_source=test")
	page := response.Body.String()
	if response.Code != 200 || response.Header().Get("X-Robots-Tag") != "" || strings.Contains(page, "noindex") {
		t.Fatal("a blog with real published posts must be indexable")
	}
	if strings.Count(page, `<article class="blog-card">`) != 2 || strings.Index(page, newer.Title) > strings.Index(page, older.Title) {
		t.Error("expected two real posts, newest first")
	}
	for _, hidden := range []string{draft.Title, future.Title, "primjer-blog-objave"} {
		if strings.Contains(page, hidden) {
			t.Errorf("listing leaked %q", hidden)
		}
	}
	if !strings.Contains(page, `href="https://eriksoftware.hr/blog"`) {
		t.Error("query string leaked into canonical URL")
	}
	for _, path := range []string{older.Path(), newer.Path()} {
		if response := request(t, h, path); response.Code != 200 || response.Header().Get("X-Robots-Tag") != "" {
			t.Errorf("published route %s failed", path)
		}
	}
	for _, path := range []string{draft.Path(), future.Path(), "/blog/missing", "/blog/missing/", newer.Path() + "/extra"} {
		if response := request(t, h, path); response.Code != 404 {
			t.Errorf("hidden/unknown route %s = %d, want 404", path, response.Code)
		}
	}
	for path, want := range map[string]string{"/blog/?x=1": "/blog?x=1", newer.Path() + "/?x=1": newer.Path() + "?x=1"} {
		r := request(t, h, path)
		if r.Code != http.StatusPermanentRedirect || r.Header().Get("Location") != want {
			t.Errorf("canonical redirect %s: %d %s", path, r.Code, r.Header().Get("Location"))
		}
	}
	graph := schemaGraph(t, page)
	list := schemaNode(t, graph, "ItemList")
	if list["numberOfItems"] != float64(2) {
		t.Error("collection schema does not match the visible listing")
	}
	schemaNode(t, graph, "CollectionPage")
	schemaNode(t, graph, "Blog")
	schemaNode(t, graph, "BreadcrumbList")
}

func TestArticleMetadataAndContent(t *testing.T) {
	p := publishedFixture("dugi-naslov", "2024-03-02T09:00:00+01:00")
	p.Title = `Dugačak naslov o automatizaciji, slovima čćđšž i sigurnosti < & " </script><script>alert(1)</script>`
	p.Description = `Opis s navodnicima " i znakom & za provjeru sigurnog kodiranja.`
	p.ModifiedAt = "2024-04-03T10:00:00+02:00"
	related := publishedFixture("povezani", "2024-02-02T09:00:00+01:00")
	h := fixtureHandler(t, p, related)
	r := request(t, h, p.Path())
	page := r.Body.String()
	for _, want := range []string{
		`<meta property="og:type" content="article">`,
		`<meta property="og:image" content="https://eriksoftware.hr/assets/blog/example/social.jpg">`,
		`<meta property="article:published_time" content="2024-03-02T09:00:00+01:00">`,
		`<meta property="article:modified_time" content="2024-04-03T10:00:00+02:00">`,
		`<meta name="robots" content="max-image-preview:large">`,
		`<link rel="canonical" href="` + p.URL() + `">`,
		`srcset="/assets/blog/example/hero-800.jpg 800w, /assets/blog/example/hero-1200.jpg 1200w, /assets/blog/example/hero.jpg 2000w"`,
		`fetchpriority="high"`, `Erik Jermaniš`, `Ažurirano`, `03/04/2024`,
		`href="/blog/povezani"`, `aria-label="Sadržaj članka"`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("article missing %q", want)
		}
	}
	if len(regexp.MustCompile(`<h1(?:\s|>)`).FindAllString(page, -1)) != 1 || strings.Contains(page, "noindex") || strings.Contains(page, `name="keywords"`) {
		t.Error("article must have one h1, be indexable, and omit meta keywords")
	}
	if strings.Contains(page, `</script><script>alert(1)</script>`) {
		t.Error("unescaped script content")
	}
	graph := schemaGraph(t, page)
	article := schemaNode(t, graph, "BlogPosting")
	if article["headline"] != p.Title || article["description"] != p.Description || article["dateModified"] != p.ModifiedAt || article["inLanguage"] != "hr" {
		t.Errorf("incorrect article schema: %+v", article)
	}
	if article["author"].(map[string]any)["@id"] != blog.BaseURL+"/#erik" || article["publisher"].(map[string]any)["@id"] != blog.BaseURL+"/#organization" {
		t.Error("author/publisher entities are disconnected")
	}
	if article["wordCount"].(float64) < 500 || article["timeRequired"] == "PT0M" {
		t.Error("reading time/word count not derived")
	}
	crumbs := schemaNode(t, graph, "BreadcrumbList")["itemListElement"].([]any)
	if len(crumbs) != 3 || crumbs[2].(map[string]any)["item"] != p.URL() {
		t.Error("article breadcrumbs incorrect")
	}
	assertLocalLinks(t, h, page)
}

func TestExampleAndDisabledStates(t *testing.T) {
	h := fixtureHandler(t, blog.ExamplePost())
	r := request(t, h, "/blog/primjer-blog-objave")
	page := r.Body.String()
	if r.Code != 200 || r.Header().Get("X-Robots-Tag") != "noindex, follow" || !strings.Contains(page, `content="noindex, follow"`) || !strings.Contains(page, "Primjer strukture, ne studija slučaja.") {
		t.Fatal("example must be clearly labelled and noindex in HTML and HTTP")
	}
	for _, want := range []string{"<table>", "<caption>", "<figure", "<figcaption>", "<blockquote>", "<pre ", `role="doc-noteref"`, `aria-labelledby="chart-title chart-description"`, `scope="col"`, `class="article-callout"`} {
		if !strings.Contains(page, want) {
			t.Errorf("showcase missing %q", want)
		}
	}
	assertLocalLinks(t, h, page)
	for _, endpoint := range []string{"/sitemap.xml", "/blog/feed.xml"} {
		if strings.Contains(request(t, h, endpoint).Body.String(), "primjer-blog-objave") {
			t.Errorf("example leaked into %s", endpoint)
		}
	}
	h = fixtureHandler(t)
	if request(t, h, "/blog/primjer-blog-objave").Code != 404 {
		t.Error("unregistered example is still reachable")
	}
	page = request(t, h, "/blog").Body.String()
	if !strings.Contains(page, "Prvi tekstovi stižu uskoro.") || !strings.Contains(page, "noindex") || strings.Contains(page, `class="blog-card"`) {
		t.Error("empty blog must have an honest noindex empty state")
	}
}

func TestSitemapAndRSSUsePublicationState(t *testing.T) {
	older := publishedFixture("stariji", "2024-01-02T09:00:00+01:00")
	older.ModifiedAt = "2024-05-03T10:00:00+02:00"
	newer := publishedFixture("noviji", "2024-03-02T09:00:00+01:00")
	newer.Title = `Provjera <naslova> & slova čćđšž`
	draft := publishedFixture("nacrt", "2024-02-02T09:00:00+01:00")
	draft.Status = blog.Draft
	future := publishedFixture("buduci", "2099-01-02T09:00:00+01:00")
	h := fixtureHandler(t, older, newer, draft, future, blog.ExamplePost())
	sm := request(t, h, "/sitemap.xml")
	if sm.Header().Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Error("incorrect sitemap MIME type")
	}
	var sitemap struct {
		XMLName xml.Name `xml:"urlset"`
		URLs    []struct {
			Location string `xml:"loc"`
			Modified string `xml:"lastmod"`
			Image    struct {
				Location string `xml:"loc"`
			} `xml:"image"`
		} `xml:"url"`
	}
	if err := xml.Unmarshal(sm.Body.Bytes(), &sitemap); err != nil {
		t.Fatal(err)
	}
	if sitemap.XMLName.Space != "http://www.sitemaps.org/schemas/sitemap/0.9" || len(sitemap.URLs) != 4 || sitemap.URLs[1].Modified != older.ModifiedAt {
		t.Errorf("incorrect sitemap: %+v", sitemap)
	}
	if sitemap.URLs[2].Location != newer.URL() || sitemap.URLs[2].Image.Location != newer.Hero.URL() {
		t.Error("sitemap missing article or image URL")
	}
	r := request(t, h, "/blog/feed.xml")
	if r.Header().Get("Content-Type") != "application/rss+xml; charset=utf-8" {
		t.Error("incorrect RSS MIME type")
	}
	var feed rss
	if err := xml.Unmarshal(r.Body.Bytes(), &feed); err != nil {
		t.Fatal(err)
	}
	if len(feed.Channel.Items) != 2 || feed.Channel.Items[0].Title != newer.Title || feed.Channel.Items[0].GUID.URL != newer.URL() || !feed.Channel.Items[0].GUID.IsPermalink {
		t.Errorf("incorrect feed items: %+v", feed.Channel.Items)
	}
	if _, err := time.Parse(time.RFC1123Z, feed.Channel.Items[0].Published); err != nil {
		t.Error(err)
	}
	if feed.Channel.LastBuildDate != "Fri, 03 May 2024 10:00:00 +0200" || !strings.Contains(r.Body.String(), `<atom:link href="https://eriksoftware.hr/blog/feed.xml" rel="self" type="application/rss+xml">`) {
		t.Error("feed must use real modification dates and advertise its canonical self URL")
	}
	for _, body := range []string{sm.Body.String(), r.Body.String()} {
		for _, hidden := range []string{draft.Slug, future.Slug, "primjer-blog-objave"} {
			if strings.Contains(body, hidden) {
				t.Errorf("XML leaked hidden post %q", hidden)
			}
		}
	}
}

func schemaGraph(t *testing.T, page string) []map[string]any {
	t.Helper()
	var graph []map[string]any
	for _, match := range regexp.MustCompile(`(?s)<script type="application/ld\+json">(.*?)</script>`).FindAllStringSubmatch(page, -1) {
		var data struct {
			Context string           `json:"@context"`
			Graph   []map[string]any `json:"@graph"`
		}
		if err := json.Unmarshal([]byte(match[1]), &data); err != nil {
			t.Fatalf("invalid structured data: %v", err)
		}
		if data.Context != "https://schema.org" {
			t.Error("missing schema context")
		}
		graph = append(graph, data.Graph...)
	}
	return graph
}

func schemaNode(t *testing.T, graph []map[string]any, kind string) map[string]any {
	t.Helper()
	for _, node := range graph {
		if node["@type"] == kind {
			return node
		}
	}
	t.Fatalf("missing schema %s", kind)
	return nil
}

func assertLocalLinks(t *testing.T, h http.Handler, page string) {
	t.Helper()
	ids := make(map[string]bool)
	for _, m := range regexp.MustCompile(`\bid="([^"]+)"`).FindAllStringSubmatch(page, -1) {
		if ids[m[1]] {
			t.Errorf("duplicate id %s", m[1])
		}
		ids[m[1]] = true
	}
	for _, m := range regexp.MustCompile(`\b(?:href|src)="([^"]+)"`).FindAllStringSubmatch(page, -1) {
		link := html.UnescapeString(m[1])
		if strings.HasPrefix(link, "#") && !ids[strings.TrimPrefix(link, "#")] {
			t.Errorf("broken fragment %s", link)
		}
		if strings.HasPrefix(link, "/") && !strings.HasPrefix(link, "//") {
			parsed, err := url.Parse(link)
			if err != nil {
				t.Fatal(err)
			}
			r := request(t, h, parsed.RequestURI())
			if r.Code != 200 {
				t.Errorf("broken internal link/asset %s (%d)", link, r.Code)
			} else if parsed.Fragment != "" && !strings.Contains(r.Body.String(), `id="`+parsed.Fragment+`"`) {
				t.Errorf("broken cross-page fragment %s", link)
			}
		}
	}
}
