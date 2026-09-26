# Erik Software

Small Croatian-language site built with Go, [templ](https://templ.guide/), and plain CSS. Go 1.22+ and [just](https://just.systems/) are the only local prerequisites. The templ command runs through the version pinned in `go.mod`; no Node or global templ installation is needed.

## Local commands

| Command | Purpose |
| --- | --- |
| `just dev` | Generate templates and start the site at `http://localhost:8080`. Set `PORT=3000 just dev` to use another port. Restart after editing templates or assets. |
| `just generate` | Regenerate `pages/*_templ.go` after editing `.templ` files. |
| `just build` | Generate templates and build `bin/eriksoftware` with static files embedded. Run it with `./bin/eriksoftware`. |
| `just test` | Generate templates and run Go tests. |

Routes: `/` (homepage), `/blog` (WIP), `/favicon.ico`, `/robots.txt`, `/sitemap.xml`. Unmatched paths return 404.

## Content and assets

The page templates are in `pages/`. Add CSS, images, and other static files to `public/` and refer to them at `/assets/<filename>`; Go embeds that directory into the binary. The navbar logo is `public/eriksoftware_logo.svg`. The shared social image is `public/social-sharing.jpg`. The SVG and Apple touch icons live in `public/`, and the ICO is also served at `/favicon.ico` for browser fallback.

## SEO checklist for the next iteration

The current layout provides per-page titles and descriptions, canonical URLs, Croatian language and social image tags, favicon links, Person + Organization + WebSite JSON-LD, robots.txt, and a sitemap. The empty blog is `noindex` and excluded from the sitemap until there are posts. Before publishing or when the relevant details exist:

1. Review the drafted homepage title and meta description in `pages/home.templ` against the final positioning and copy. Update `/blog` metadata when it has real content.
2. `https://eriksoftware.hr/` (without `www`) is the preferred public URL. Keep canonical tags, JSON-LD, sitemap and robots.txt consistent if it changes.
3. If public business details such as a service area or registered legal name are provided later, add only verified, site-visible information to the Organization schema. Erik Jermaniš's LinkedIn is associated with the Person, not the business.
4. When publishing posts, add them to `public/sitemap.xml`, give each a unique canonical/title/description and suitable BlogPosting JSON-LD, and remove the blog index's `noindex`. Submit the sitemap in Google Search Console after launch.
