package pages

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ErikJermanis/eriksoftware.hr/blog"
	"github.com/a-h/templ"
)

type metadata struct {
	Title, Description, Canonical string
	NoIndex                       bool
	Kind                          string
	Social                        blog.Image
	PublishedAt, ModifiedAt       string
	Schema                        any
}

func (m metadata) kind() string {
	if m.Kind != "" {
		return m.Kind
	}
	return "website"
}

func (m metadata) social() blog.Image {
	if m.Social.Src != "" {
		return m.Social
	}
	return blog.Image{Src: "/assets/social-sharing.jpg", Alt: "Erik Software — software i automatizacija za poduzetnike", Width: 1200, Height: 630}
}

// json.Marshal escapes <, > and &, so even a title containing </script> cannot
// break out of the JSON-LD element. Never build JSON-LD by string interpolation.
func structuredData(data any) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, `<script type="application/ld+json">%s</script>`, encoded)
		return err
	})
}

func breadcrumbs(post *blog.Post) map[string]any {
	items := []map[string]any{
		{"@type": "ListItem", "position": 1, "name": "Naslovnica", "item": blog.BaseURL + "/"},
		{"@type": "ListItem", "position": 2, "name": "Blog", "item": blog.BaseURL + "/blog"},
	}
	url := blog.BaseURL + "/blog"
	if post != nil {
		url = post.URL()
		items = append(items, map[string]any{"@type": "ListItem", "position": 3, "name": post.Title, "item": post.URL()})
	}
	return map[string]any{"@type": "BreadcrumbList", "@id": url + "#breadcrumbs", "itemListElement": items}
}

func collectionSchema(posts []blog.Post) map[string]any {
	items := make([]map[string]any, 0, len(posts))
	for i, post := range posts {
		items = append(items, map[string]any{"@type": "ListItem", "position": i + 1, "name": post.Title, "url": post.URL()})
	}
	url := blog.BaseURL + "/blog"
	return map[string]any{"@context": "https://schema.org", "@graph": []any{
		map[string]any{
			"@type": "CollectionPage", "@id": url + "#webpage", "url": url,
			"name": blogTitle, "description": blogDescription, "inLanguage": "hr",
			"isPartOf":   map[string]string{"@id": blog.BaseURL + "/#website"},
			"mainEntity": map[string]string{"@id": url + "#posts"},
			"breadcrumb": map[string]string{"@id": url + "#breadcrumbs"},
		},
		map[string]any{
			"@type": "Blog", "@id": url + "#blog", "url": url,
			"name": "Erik Software — Blog", "inLanguage": "hr",
			"publisher": map[string]string{"@id": blog.BaseURL + "/#organization"},
			"author":    map[string]string{"@id": blog.BaseURL + "/#erik"},
		},
		map[string]any{
			"@type": "ItemList", "@id": url + "#posts",
			"numberOfItems": len(items), "itemListElement": items,
		},
		breadcrumbs(nil),
	}}
}

func articleSchema(post blog.Post) map[string]any {
	article := map[string]any{
		"@type": "BlogPosting", "@id": post.URL() + "#article", "url": post.URL(),
		"headline": post.Title, "description": post.Description, "inLanguage": "hr",
		"datePublished": post.PublishedAt, "wordCount": post.WordCount,
		"timeRequired":     fmt.Sprintf("PT%dM", post.ReadingMinutes),
		"image":            []string{post.Social.URL(), post.Hero.URL()},
		"mainEntityOfPage": map[string]string{"@id": post.URL() + "#webpage"},
		"isPartOf":         map[string]string{"@id": blog.BaseURL + "/blog#blog"},
		"author":           map[string]string{"@id": blog.BaseURL + "/#erik"},
		"publisher":        map[string]string{"@id": blog.BaseURL + "/#organization"},
	}
	if post.ModifiedAt != "" {
		article["dateModified"] = post.ModifiedAt
	}
	if len(post.Topics) > 0 {
		article["articleSection"] = post.Topics
	}
	return map[string]any{"@context": "https://schema.org", "@graph": []any{
		map[string]any{
			"@type": "WebPage", "@id": post.URL() + "#webpage", "url": post.URL(),
			"name": post.Title, "inLanguage": "hr",
			"isPartOf":   map[string]string{"@id": blog.BaseURL + "/#website"},
			"breadcrumb": map[string]string{"@id": post.URL() + "#breadcrumbs"},
			"primaryImageOfPage": map[string]any{
				"@type": "ImageObject", "url": post.Social.URL(),
				"width": post.Social.Width, "height": post.Social.Height,
			},
		},
		article,
		breadcrumbs(&post),
	}}
}
