package blog

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/a-h/templ"
)

func TestMetadataValidation(t *testing.T) {
	for name, mutate := range map[string]func(*Post){
		"unsafe slug":          func(p *Post) { p.Slug = "../etc/passwd" },
		"missing title":        func(p *Post) { p.Title = " " },
		"missing description":  func(p *Post) { p.Description = "" },
		"missing body":         func(p *Post) { p.Body = nil },
		"missing status":       func(p *Post) { p.Status = "" },
		"bad date":             func(p *Post) { p.PublishedAt = "2026-02-30" },
		"bad modified date":    func(p *Post) { p.ModifiedAt = "yesterday" },
		"modified before date": func(p *Post) { p.ModifiedAt = "2024-01-02T09:00:00+01:00" },
		"wrong hero ratio":     func(p *Post) { p.Hero.Height = 630 },
		"wrong file size":      func(p *Post) { p.Hero.Width, p.Hero.Height = 1000, 400 },
		"missing alt":          func(p *Post) { p.Hero.Alt = "" },
		"missing hero":         func(p *Post) { p.Hero.Src = "/assets/missing.jpg" },
		"external hero":        func(p *Post) { p.Hero.Src = "https://example.com/image.jpg" },
		"asset traversal":      func(p *Post) { p.Hero.Src = "/assets/../hero.jpg" },
		"wrong extension":      func(p *Post) { p.Hero.Src = "/assets/blog/example/hero.gif" },
		"SVG social preview":   func(p *Post) { p.Social.Src = "/assets/blog/example/hero.svg" },
		"tiny social preview":  func(p *Post) { p.Social.Width = 800 },
		"bad heading":          func(p *Post) { p.Headings = []Heading{{ID: "missing", Title: "Missing"}} },
		"repeated heading":     func(p *Post) { p.Headings = []Heading{exampleHeadings[0], exampleHeadings[0]} },
		"bad variant ratio":    func(p *Post) { p.Hero.Variants[0].Height = 100 },
		"missing variant":      func(p *Post) { p.Hero.Variants[0].Src = "/assets/missing.jpg" },
		"duplicate variant":    func(p *Post) { p.Hero.Variants = []Source{p.Hero.Variants[0], p.Hero.Variants[0]} },
		"larger variant":       func(p *Post) { p.Hero.Variants[0].Width, p.Hero.Variants[0].Height = 2500, 1000 },
		"body render error": func(p *Post) {
			p.Body = templ.ComponentFunc(func(context.Context, io.Writer) error { return errors.New("broken") })
		},
	} {
		t.Run(name, func(t *testing.T) {
			p := ExamplePost()
			mutate(&p)
			if _, err := NewCatalog([]Post{p}, os.DirFS("../public")); err == nil {
				t.Error("invalid post accepted")
			}
		})
	}
	if _, err := NewCatalog([]Post{ExamplePost(), ExamplePost()}, os.DirFS("../public")); err == nil {
		t.Error("duplicate slug accepted")
	}
}

func TestCatalogPublicationAndRelatedPosts(t *testing.T) {
	now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	older := ExamplePost()
	older.Slug, older.Status, older.PublishedAt = "stariji", Published, "2024-01-01T09:00:00+01:00"
	older.Excerpt = ""
	newer := older
	newer.Slug, newer.PublishedAt = "noviji", "2024-02-01T09:00:00+01:00"
	unrelated := older
	unrelated.Slug, unrelated.Topics = "druga-tema", []string{"Drugo"}
	draft := older
	draft.Slug, draft.Status = "nacrt", Draft
	future := older
	future.Slug, future.PublishedAt = "buduci", "2099-01-01T09:00:00+01:00"
	c, err := NewCatalog([]Post{older, newer, unrelated, draft, future, ExamplePost()}, os.DirFS("../public"))
	if err != nil {
		t.Fatal(err)
	}
	if list := c.Listed(now); len(list) != 3 || list[0].Slug != newer.Slug {
		t.Fatalf("incorrect listing: %+v", list)
	}
	p, ok := c.Lookup(older.Slug, now)
	if !ok || p.Excerpt != p.Description || p.ReadingMinutes < 2 || p.WordCount < 500 || p.Date() != "01/01/2024" {
		t.Error("derived article fields incorrect")
	}
	if related := c.Related(p, now); len(related) != 1 || related[0].Slug != newer.Slug {
		t.Errorf("related articles must exclude self, examples, drafts, future and unrelated posts: %+v", related)
	}
	for _, slug := range []string{draft.Slug, future.Slug, "missing"} {
		if _, ok := c.Lookup(slug, now); ok {
			t.Errorf("unexpectedly visible: %s", slug)
		}
	}
	if _, ok := c.Lookup(future.Slug, time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)); !ok {
		t.Error("future publication did not become available when its date passed")
	}
	if !strings.Contains(p.Hero.Srcset(), "1200w") || p.Hero.MIME() != "image/jpeg" || p.Hero.URL() != BaseURL+p.Hero.Src {
		t.Error("image helpers incorrect")
	}
}

func TestEmptyCatalogAndExamples(t *testing.T) {
	now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	c, err := NewCatalog(nil, os.DirFS("../public"))
	if err != nil || len(c.Listed(now)) != 0 || len(c.Published(now)) != 0 {
		t.Fatal("empty registry must work without changing templates or tests")
	}
	c, err = NewCatalog([]Post{ExamplePost()}, os.DirFS("../public"))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Listed(now)) != 1 || len(c.Published(now)) != 0 {
		t.Error("example should be listed, never published")
	}
	p, ok := c.Lookup("primjer-blog-objave", now)
	if !ok || p.LastModified() != p.Published() {
		t.Error("unmodified example must use its publication date")
	}
}
