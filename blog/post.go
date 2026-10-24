// Package blog holds trusted, compiled article content and its publication metadata.
package blog

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

const BaseURL = "https://eriksoftware.hr"

type Status string

const (
	Published Status = "published"
	Draft     Status = "draft"
	Example   Status = "example"
)

type Source struct {
	Src           string
	Width, Height int
}

type Image struct {
	Src           string
	Alt           string
	Width, Height int
	Variants      []Source
}

func (i Image) URL() string { return BaseURL + i.Src }

func (i Image) Srcset() string {
	if len(i.Variants) == 0 {
		return ""
	}
	sources := make([]string, 0, len(i.Variants)+1)
	for _, source := range i.Variants {
		sources = append(sources, fmt.Sprintf("%s %dw", source.Src, source.Width))
	}
	sources = append(sources, fmt.Sprintf("%s %dw", i.Src, i.Width))
	return strings.Join(sources, ", ")
}

func (i Image) MIME() string {
	switch strings.ToLower(path.Ext(i.Src)) {
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	default:
		return "image/jpeg"
	}
}

type Heading struct {
	ID, Title string
}

type Post struct {
	Slug, Title, Description string
	Excerpt                  string // Defaults to Description.
	PublishedAt, ModifiedAt  string // RFC 3339, with the actual local UTC offset.
	Status                   Status
	Topics                   []string
	Hero, Social             Image
	Headings                 []Heading
	Body                     templ.Component
	ReadingMinutes           int // Derived from rendered content, not entered by authors.
	WordCount                int
	published, modified      time.Time
}

func (p Post) Path() string         { return "/blog/" + p.Slug }
func (p Post) URL() string          { return BaseURL + p.Path() }
func (p Post) Date() string         { return p.published.Format("02/01/2006") }
func (p Post) UpdatedDate() string  { return p.modified.Format("02/01/2006") }
func (p Post) Published() time.Time { return p.published }
func (p Post) LastModified() time.Time {
	if !p.modified.IsZero() {
		return p.modified
	}
	return p.published
}

type Catalog struct {
	posts  []Post
	bySlug map[string]Post
}

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	tagsPattern = regexp.MustCompile(`<[^>]*>`)
)

// NewCatalog validates metadata/assets and estimates reading time once at startup.
// Bodies are trusted templ components, never raw Markdown or request input.
func NewCatalog(posts []Post, assets fs.FS) (*Catalog, error) {
	c := &Catalog{bySlug: make(map[string]Post)}
	for _, p := range posts {
		if !slugPattern.MatchString(p.Slug) || strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.Description) == "" || p.Body == nil {
			return nil, fmt.Errorf("post %q: slug, title, description and body are required", p.Slug)
		}
		if _, exists := c.bySlug[p.Slug]; exists {
			return nil, fmt.Errorf("duplicate post slug %q", p.Slug)
		}
		if p.Status != Published && p.Status != Draft && p.Status != Example {
			return nil, fmt.Errorf("post %q: explicitly choose published, draft or example", p.Slug)
		}
		var err error
		p.published, err = time.Parse(time.RFC3339, p.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("post %q: invalid publication timestamp: %w", p.Slug, err)
		}
		if p.ModifiedAt != "" {
			p.modified, err = time.Parse(time.RFC3339, p.ModifiedAt)
			if err != nil || p.modified.Before(p.published) {
				return nil, fmt.Errorf("post %q: modification timestamp must be valid and not precede publication", p.Slug)
			}
		}
		if p.Hero.Width*2 != p.Hero.Height*5 {
			return nil, fmt.Errorf("post %q: hero must have a 5:2 aspect ratio", p.Slug)
		}
		if err := validateImage(p.Hero, assets, false); err != nil {
			return nil, fmt.Errorf("post %q hero: %w", p.Slug, err)
		}
		if err := validateImage(p.Social, assets, true); err != nil {
			return nil, fmt.Errorf("post %q social image: %w", p.Slug, err)
		}
		widths := map[int]bool{p.Hero.Width: true}
		for _, source := range p.Hero.Variants {
			if widths[source.Width] || source.Width > p.Hero.Width {
				return nil, fmt.Errorf("post %q: hero variant widths must be unique and smaller than the original", p.Slug)
			}
			widths[source.Width] = true
			if source.Width*2 != source.Height*5 {
				return nil, fmt.Errorf("post %q: hero variants must also be 5:2", p.Slug)
			}
			if err := validateImage(Image{Src: source.Src, Width: source.Width, Height: source.Height, Alt: p.Hero.Alt}, assets, false); err != nil {
				return nil, fmt.Errorf("post %q hero variant: %w", p.Slug, err)
			}
		}
		var body bytes.Buffer
		if err := p.Body.Render(context.Background(), &body); err != nil {
			return nil, fmt.Errorf("post %q: render body: %w", p.Slug, err)
		}
		seen := make(map[string]bool)
		for _, heading := range p.Headings {
			if !slugPattern.MatchString(heading.ID) || strings.TrimSpace(heading.Title) == "" || seen[heading.ID] || !strings.Contains(body.String(), `id="`+heading.ID+`"`) {
				return nil, fmt.Errorf("post %q: invalid, duplicate or missing section %q", p.Slug, heading.ID)
			}
			seen[heading.ID] = true
		}
		// An estimate only: strip trusted rendered markup, then decode entities.
		p.WordCount = len(strings.Fields(html.UnescapeString(tagsPattern.ReplaceAllString(body.String(), " "))))
		p.ReadingMinutes = max(1, (p.WordCount+199)/200)
		if p.Excerpt == "" {
			p.Excerpt = p.Description
		}
		c.posts = append(c.posts, p)
		c.bySlug[p.Slug] = p
	}
	sort.Slice(c.posts, func(i, j int) bool {
		if c.posts[i].published.Equal(c.posts[j].published) {
			return c.posts[i].Slug < c.posts[j].Slug
		}
		return c.posts[i].published.After(c.posts[j].published)
	})
	return c, nil
}

func validateImage(img Image, assets fs.FS, social bool) error {
	if !strings.HasPrefix(img.Src, "/assets/") || path.Clean(img.Src) != img.Src || img.Width <= 0 || img.Height <= 0 || strings.TrimSpace(img.Alt) == "" {
		return fmt.Errorf("local /assets/ path, dimensions and alt text are required")
	}
	ext := strings.ToLower(path.Ext(img.Src))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && (social || ext != ".svg") {
		return fmt.Errorf("use JPEG/PNG, or SVG for non-social images")
	}
	if social && (img.Width < 1200 || img.Height < 630) {
		return fmt.Errorf("social preview must be at least 1200 × 630")
	}
	f, err := assets.Open(strings.TrimPrefix(img.Src, "/assets/"))
	if err != nil {
		return err
	}
	defer f.Close()
	var width, height int
	if ext == ".svg" {
		var svg struct {
			XMLName xml.Name `xml:"svg"`
			Width   string   `xml:"width,attr"`
			Height  string   `xml:"height,attr"`
		}
		if err := xml.NewDecoder(f).Decode(&svg); err != nil {
			return err
		}
		width, _ = strconv.Atoi(svg.Width)
		height, _ = strconv.Atoi(svg.Height)
	} else {
		config, format, err := image.DecodeConfig(f)
		if err != nil {
			return err
		}
		if (format == "jpeg" && ext != ".jpg" && ext != ".jpeg") || (format == "png" && ext != ".png") {
			return fmt.Errorf("image extension does not match its format")
		}
		width, height = config.Width, config.Height
	}
	if width != img.Width || height != img.Height {
		return fmt.Errorf("declared %d × %d, file is %d × %d", img.Width, img.Height, width, height)
	}
	return nil
}

func (c *Catalog) Lookup(slug string, now time.Time) (Post, bool) {
	p, ok := c.bySlug[slug]
	return p, ok && p.Status != Draft && !p.published.After(now)
}

func (c *Catalog) Published(now time.Time) []Post {
	var result []Post
	for _, p := range c.posts {
		if p.Status == Published && !p.published.After(now) {
			result = append(result, p)
		}
	}
	return result
}

// Listed shows examples only until the first real article is published.
func (c *Catalog) Listed(now time.Time) []Post {
	if posts := c.Published(now); len(posts) > 0 {
		return posts
	}
	var result []Post
	for _, p := range c.posts {
		if p.Status == Example && !p.published.After(now) {
			result = append(result, p)
		}
	}
	return result
}

func (c *Catalog) Related(post Post, now time.Time) []Post {
	var related []Post
	for _, candidate := range c.Published(now) {
		if candidate.Slug == post.Slug {
			continue
		}
		matches := false
		for _, a := range post.Topics {
			for _, b := range candidate.Topics {
				matches = matches || strings.EqualFold(a, b)
			}
		}
		if matches {
			related = append(related, candidate)
			if len(related) == 3 {
				break
			}
		}
	}
	return related
}
