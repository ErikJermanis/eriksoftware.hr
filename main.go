package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/ErikJermanis/eriksoftware.hr/pages"
	"github.com/a-h/templ"
)

//go:embed public
var publicFiles embed.FS

func newHandler() (http.Handler, error) {
	static, err := fs.Sub(publicFiles, "public")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", templ.Handler(pages.Home()))
	mux.Handle("GET /blog", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Robots-Tag", "noindex, follow")
		templ.Handler(pages.Blog()).ServeHTTP(w, r)
	}))
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(static))))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "favicon.ico")
	})
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "robots.txt")
	})
	mux.HandleFunc("GET /sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, static, "sitemap.xml")
	})

	return mux, nil
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

	log.Printf("Listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
