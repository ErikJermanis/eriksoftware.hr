package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/ErikJermanis/eriksoftware.hr/blog"
)

type sitemap struct {
	XMLName    xml.Name     `xml:"urlset"`
	XMLNS      string       `xml:"xmlns,attr"`
	XMLNSImage string       `xml:"xmlns:image,attr"`
	URLs       []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Location string        `xml:"loc"`
	Modified string        `xml:"lastmod,omitempty"`
	Image    *sitemapImage `xml:"image:image,omitempty"`
}

type sitemapImage struct {
	Location string `xml:"image:loc"`
}

func serveSitemap(w http.ResponseWriter, posts []blog.Post) {
	sm := sitemap{
		XMLNS:      "http://www.sitemaps.org/schemas/sitemap/0.9",
		XMLNSImage: "http://www.google.com/schemas/sitemap-image/1.1",
		URLs:       []sitemapURL{{Location: blog.BaseURL + "/"}},
	}
	if len(posts) > 0 {
		latest := latestModification(posts)
		sm.URLs = append(sm.URLs, sitemapURL{Location: blog.BaseURL + "/blog", Modified: latest.Format(time.RFC3339)})
		for _, post := range posts {
			sm.URLs = append(sm.URLs, sitemapURL{Location: post.URL(), Modified: post.LastModified().Format(time.RFC3339), Image: &sitemapImage{Location: post.Hero.URL()}})
		}
	}
	serveXML(w, "application/xml; charset=utf-8", sm)
}

type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	AtomNS  string     `xml:"xmlns:atom,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	LastBuildDate string    `xml:"lastBuildDate,omitempty"`
	Self          rssSelf   `xml:"atom:link"`
	Items         []rssItem `xml:"item"`
}

type rssSelf struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type rssItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	GUID        rssGUID  `xml:"guid"`
	Description string   `xml:"description"`
	Author      string   `xml:"author"`
	Published   string   `xml:"pubDate"`
	Categories  []string `xml:"category"`
}

type rssGUID struct {
	IsPermalink bool   `xml:"isPermaLink,attr"`
	URL         string `xml:",chardata"`
}

func serveRSS(w http.ResponseWriter, posts []blog.Post) {
	feed := rss{Version: "2.0", AtomNS: "http://www.w3.org/2005/Atom", Channel: rssChannel{
		Title:       "Erik Software — Blog",
		Link:        blog.BaseURL + "/blog",
		Description: "O softveru, automatizaciji poslovanja i odgovornom korištenju AI alata.",
		Language:    "hr",
		Self:        rssSelf{Href: blog.BaseURL + "/blog/feed.xml", Rel: "self", Type: "application/rss+xml"},
	}}
	if len(posts) > 0 {
		feed.Channel.LastBuildDate = latestModification(posts).Format(time.RFC1123Z)
	}
	for _, post := range posts {
		feed.Channel.Items = append(feed.Channel.Items, rssItem{
			Title: post.Title, Link: post.URL(), GUID: rssGUID{IsPermalink: true, URL: post.URL()},
			Description: post.Excerpt, Author: "erik@eriksoftware.hr (Erik Jermaniš)",
			Published: post.Published().Format(time.RFC1123Z), Categories: post.Topics,
		})
	}
	w.Header().Set("X-Robots-Tag", "noindex")
	serveXML(w, "application/rss+xml; charset=utf-8", feed)
}

func latestModification(posts []blog.Post) time.Time {
	var latest time.Time
	for _, post := range posts {
		if modified := post.LastModified(); modified.After(latest) {
			latest = modified
		}
	}
	return latest
}

func serveXML(w http.ResponseWriter, contentType string, value any) {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		log.Printf("encode XML: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(append([]byte(xml.Header), data...))
}
