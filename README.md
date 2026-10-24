# Erik Software

Croatian-language site built with Go, [templ](https://templ.guide/), and plain CSS. Requires Go 1.22+ and [just](https://just.systems/). No Node, database, or global templ installation; the templ version is pinned in `go.mod`.

| Command | Purpose |
| --- | --- |
| `just dev` | Generate templates and serve `http://localhost:8080`. Override with `PORT=3000`; restart after edits. |
| `just generate` | Regenerate all `*_templ.go` files after editing `.templ` sources. Never hand-edit generated files. |
| `just test` | Generate templates and run Go tests. |
| `just build` | Build `bin/eriksoftware`, including all content and assets. Run with `./bin/eriksoftware`. |
| `just deploy` | Build for Linux, upload to the configured server, and restart its service. Requires SSH access. |

## Content

- `pages/`: site layout, homepage, blog listing, and shared article page.
- `blog/`: compiled post bodies, metadata, and reusable article components. **Read [the authoring guide](blog/README.md) before adding a post.**
- `public/`: embedded files served at `/assets/<path>`. Keep post images in `public/blog/<slug>/`.

Routes: `/`, `/blog`, `/blog/{slug}`, `/blog/feed.xml`, `/sitemap.xml`, `/robots.txt`, `/favicon.ico`. Unknown or unpublished posts return 404; trailing blog slashes redirect to the canonical path.

## Publication and SEO

`blog/posts.go` is the only publication registry. Cards, related posts, structured data, sitemap, and RSS derive from it. Drafts and future-dated posts stay hidden. The example is disabled and retained in `blog/` only as an authoring reference.

The preferred origin is `https://eriksoftware.hr` (no `www`). Real posts receive unique canonical/social metadata and BlogPosting + breadcrumb JSON-LD; the listing becomes indexable automatically with the first real post. Shared Person, Organization, and WebSite schema stays in `pages/layout.templ`; add only verified, visible business details.

After publishing, submit `/sitemap.xml` in Google Search Console and check a real article with Google's Rich Results Test and PageSpeed Insights. Do not manually maintain sitemap entries or invent modification dates.
