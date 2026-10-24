package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ErikJermanis/eriksoftware.hr/blog"
	"github.com/ErikJermanis/eriksoftware.hr/pages"
	"github.com/a-h/templ"
)

//go:embed public
var publicFiles embed.FS

func newHandler() (http.Handler, error) {
	return handlerWithPosts(blog.Posts())
}

func handlerWithPosts(posts []blog.Post) (http.Handler, error) {
	static, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return nil, err
	}
	catalog, err := blog.NewCatalog(posts, static)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(pages.Home()))
	mux.Handle("GET /blog", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		noIndex := len(catalog.Published(now)) == 0
		if noIndex {
			w.Header().Set("X-Robots-Tag", "noindex, follow")
		}
		templ.Handler(pages.Blog(catalog.Listed(now), noIndex)).ServeHTTP(w, r)
	}))
	mux.HandleFunc("GET /blog/{$}", func(w http.ResponseWriter, r *http.Request) {
		redirectCanonical(w, r, "/blog")
	})
	mux.HandleFunc("GET /blog/feed.xml", func(w http.ResponseWriter, r *http.Request) {
		serveRSS(w, catalog.Published(time.Now()))
	})
	mux.HandleFunc("GET /blog/{slug}", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		post, ok := catalog.Lookup(r.PathValue("slug"), now)
		if !ok {
			http.NotFound(w, r)
			return
		}
		if post.Status != blog.Published {
			w.Header().Set("X-Robots-Tag", "noindex, follow")
		}
		templ.Handler(pages.BlogPost(post, catalog.Related(post, now))).ServeHTTP(w, r)
	})
	mux.HandleFunc("GET /blog/{slug}/{$}", func(w http.ResponseWriter, r *http.Request) {
		post, ok := catalog.Lookup(r.PathValue("slug"), time.Now())
		if !ok {
			http.NotFound(w, r)
			return
		}
		redirectCanonical(w, r, post.Path())
	})
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(static))))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "favicon.ico")
	})
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "robots.txt")
	})
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		serveSitemap(w, catalog.Published(time.Now()))
	})

	return mux, nil
}

func redirectCanonical(w http.ResponseWriter, r *http.Request, path string) {
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, path, http.StatusPermanentRedirect)
}

func main() {
	handler, err := newHandler()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := "127.0.0.1:" + port
	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
