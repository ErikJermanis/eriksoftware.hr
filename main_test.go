package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomepageMetadataAndSchema(t *testing.T) {
	handler := testHandler(t)
	response := request(t, handler, "/")
	if response.Code != http.StatusOK {
		t.Fatalf("homepage status = %d", response.Code)
	}

	page := response.Body.String()
	for _, want := range []string{
		`<html lang="hr">`,
		`<link rel="canonical" href="https://eriksoftware.hr/">`,
		`<meta name="description" content=`,
		`<meta property="og:locale" content="hr_HR">`,
		`<meta property="og:image" content="https://eriksoftware.hr/assets/social-sharing.jpg">`,
		`<meta property="og:image:width" content="1200">`,
		`<meta property="og:image:height" content="630">`,
		`<meta name="twitter:card" content="summary_large_image">`,
		`<meta name="twitter:image" content="https://eriksoftware.hr/assets/social-sharing.jpg">`,
		`<link rel="icon" href="/favicon.ico">`,
		`<link rel="icon" type="image/svg+xml" sizes="any" href="/assets/favicon.svg">`,
		`<link rel="apple-touch-icon" sizes="180x180" href="/assets/apple-touch-icon.png">`,
		`href="mailto:erik@eriksoftware.hr"`,
		`<h1 id="home-title">`,
		`id="kako-radim"`,
		`id="konzultacije"`,
		`/assets/eriksoftware_logo.svg`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("homepage missing %q", want)
		}
	}
	if strings.Contains(page, `content="noindex`) {
		t.Error("homepage unexpectedly noindex")
	}

	_, script, ok := strings.Cut(page, `<script type="application/ld+json">`)
	if !ok {
		t.Fatal("JSON-LD script missing")
	}
	script, _, ok = strings.Cut(script, `</script>`)
	if !ok {
		t.Fatal("JSON-LD script not closed")
	}
	var schema struct {
		Context string `json:"@context"`
		Graph   []struct {
			Type      string `json:"@type"`
			ID        string `json:"@id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			SameAs    string `json:"sameAs"`
			WorkFor   struct{ ID string `json:"@id"` } `json:"worksFor"`
			Owner     struct{ ID string `json:"@id"` } `json:"owner"`
			Employee  struct{ ID string `json:"@id"` } `json:"employee"`
			Publisher struct{ ID string `json:"@id"` } `json:"publisher"`
		} `json:"@graph"`
	}
	if err := json.Unmarshal([]byte(script), &schema); err != nil {
		t.Fatalf("invalid JSON-LD: %v", err)
	}
	if schema.Context != "https://schema.org" || len(schema.Graph) != 3 || schema.Graph[0].Type != "Person" || schema.Graph[1].Type != "Organization" || schema.Graph[2].Type != "WebSite" {
		t.Fatalf("unexpected JSON-LD: %+v", schema)
	}
	person, organization, website := schema.Graph[0], schema.Graph[1], schema.Graph[2]
	if person.SameAs != "https://www.linkedin.com/in/erik-jermanis/" || person.WorkFor.ID != organization.ID || organization.Name != "Erik Software" || organization.Email != "erik@eriksoftware.hr" || organization.Owner.ID != person.ID || organization.Employee.ID != person.ID || website.Publisher.ID != organization.ID {
		t.Fatalf("incorrect entity relationships: %+v", schema.Graph)
	}
}

func TestBlogPlaceholdersAreNotIndexed(t *testing.T) {
	handler := testHandler(t)
	response := request(t, handler, "/blog")
	if response.Code != http.StatusOK {
		t.Fatalf("blog status = %d", response.Code)
	}
	if response.Header().Get("X-Robots-Tag") != "noindex, follow" {
		t.Error("blog missing X-Robots-Tag")
	}
	page := response.Body.String()
	for _, want := range []string{
		`content="noindex, follow"`,
		`href="https://eriksoftware.hr/blog"`,
		`<h1 id="blog-title">Blog</h1>`,
		`Erik Jermaniš`,
		`datetime="2026-09-01">01/09/2026</time>`,
		`datetime="2026-09-15">15/09/2026</time>`,
		`/assets/blog-process.svg`,
		`/assets/blog-decision.svg`,
	} {
		if !strings.Contains(page, want) {
			t.Errorf("blog missing %q", want)
		}
	}
	if count := strings.Count(page, `<article class="blog-card">`); count != 2 {
		t.Errorf("blog cards = %d, want 2", count)
	}
	if strings.Contains(request(t, handler, "/sitemap.xml").Body.String(), "/blog") {
		t.Error("unfinished blog listed in sitemap")
	}
}

func TestStaticAssetsAndUnknownRoutes(t *testing.T) {
	handler := testHandler(t)
	for _, path := range []string{"/assets/eriksoftware_logo.svg", "/assets/styles.css", "/assets/social-sharing.jpg", "/assets/favicon.svg", "/assets/apple-touch-icon.png", "/assets/blog-process.svg", "/assets/blog-decision.svg", "/favicon.ico", "/robots.txt", "/sitemap.xml"} {
		if status := request(t, handler, path).Code; status != http.StatusOK {
			t.Errorf("%s status = %d", path, status)
		}
	}
	for _, path := range []string{"/missing", "/blog/post"} {
		if status := request(t, handler, path).Code; status != http.StatusNotFound {
			t.Errorf("%s status = %d", path, status)
		}
	}
}

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	handler, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func request(t *testing.T, handler http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response
}
