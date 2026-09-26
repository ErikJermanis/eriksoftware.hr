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
			Type string `json:"@type"`
		} `json:"@graph"`
	}
	if err := json.Unmarshal([]byte(script), &schema); err != nil {
		t.Fatalf("invalid JSON-LD: %v", err)
	}
	if schema.Context != "https://schema.org" || len(schema.Graph) != 2 || schema.Graph[0].Type != "Person" || schema.Graph[1].Type != "WebSite" {
		t.Fatalf("unexpected JSON-LD: %+v", schema)
	}
}

func TestBlogPlaceholderIsNotIndexed(t *testing.T) {
	handler := testHandler(t)
	response := request(t, handler, "/blog")
	if response.Code != http.StatusOK {
		t.Fatalf("blog status = %d", response.Code)
	}
	if response.Header().Get("X-Robots-Tag") != "noindex, follow" {
		t.Error("blog missing X-Robots-Tag")
	}
	for _, want := range []string{`content="noindex, follow"`, `href="https://eriksoftware.hr/blog"`, `WIP — uskoro više.`} {
		if !strings.Contains(response.Body.String(), want) {
			t.Errorf("blog missing %q", want)
		}
	}
	if strings.Contains(request(t, handler, "/sitemap.xml").Body.String(), "/blog") {
		t.Error("unfinished blog listed in sitemap")
	}
}

func TestStaticAssetsAndUnknownRoutes(t *testing.T) {
	handler := testHandler(t)
	for _, path := range []string{"/assets/eriksoftware_logo.svg", "/assets/styles.css", "/robots.txt", "/sitemap.xml"} {
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
