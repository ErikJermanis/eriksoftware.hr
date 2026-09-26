# Erik Software

Small Croatian-language site built with Go, [templ](https://templ.guide/), and plain CSS. Go 1.22+ and [just](https://just.systems/) are the only local prerequisites. The templ command runs through the version pinned in `go.mod`; no Node or global templ installation is needed.

## Local commands

| Command | Purpose |
| --- | --- |
| `just dev` | Generate templates and start the site at `http://localhost:8080`. Set `PORT=3000 just dev` to use another port. Restart after editing templates or assets. |
| `just generate` | Regenerate `pages/*_templ.go` after editing `.templ` files. |
| `just build` | Generate templates and build `bin/eriksoftware` with static files embedded. Run it with `./bin/eriksoftware`. |
| `just test` | Generate templates and run Go tests. |

Routes: `/` (homepage), `/blog` (WIP), `/robots.txt`, `/sitemap.xml`. Unmatched paths return 404.

## Content and assets

The page templates are in `pages/`. Add CSS, images, and other static files to `public/` and refer to them at `/assets/<filename>`; Go embeds that directory into the binary. The navbar logo is `public/eriksoftware_logo.svg`. When ready, add a favicon under `public/` and its `<link rel="icon">` tag in `pages/layout.templ`.

## SEO checklist for the next iteration

The current layout provides per-page titles and descriptions, canonical URLs, Croatian language and social tags, Person + WebSite JSON-LD, robots.txt, and a sitemap. The empty blog is `noindex` and excluded from the sitemap until there are posts. Before publishing or when the relevant assets/details exist:

1. Review the drafted homepage title and meta description in `pages/home.templ` against the final positioning and copy. Update `/blog` metadata when it has real content.
2. Confirm that `https://eriksoftware.hr/` (without `www`) is the preferred public URL. Update canonical tags, JSON-LD, sitemap and robots.txt together if the domain format changes.
3. Create a real social-sharing image (ideally 1200 × 630 px), place it in `public/`, then add absolute `og:image` and `twitter:image` URLs and appropriate image alt text in `pages/layout.templ`. Consider switching the Twitter card to `summary_large_image` then.
4. Add favicon assets in `public/` and the corresponding icon link(s) to the layout. The current SVG logo can be used as the starting point.
5. Provide verified public business details (official business name, service area/location, contact method and any public profile URLs). Then extend the JSON-LD with appropriate `ProfessionalService` details and approved `sameAs` links; add a contact link to the site when one is available.
6. When publishing posts, add them to `public/sitemap.xml`, give each a unique canonical/title/description and suitable BlogPosting JSON-LD, and remove the blog index's `noindex`. Submit the sitemap in Google Search Console after launch.
