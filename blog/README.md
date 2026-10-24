# Writing a post

Write in Obsidian; give an agent the `.md` file and its attachments. The agent converts it into **trusted templ content**, not a runtime Markdown import. Preserve the author's meaning, language, and voice; ask about missing assets or facts rather than inventing them. Never edit the original vault file unless asked.

Reference: [`example.md`](example.md) → [`example.templ`](example.templ), with metadata in [`posts.go`](posts.go). The Markdown is an input specimen, not a second runtime source of truth. Shared structure/styles live in `pages/blog_post.templ`, `components.templ`, and `public/styles.css`.

## Conversion checklist

1. Create `blog/<slug>.templ` (`package blog`) with the article body only. Use one unique body function. **No H1 or duplicated hero/byline**; the shared page supplies those.
2. Create its `Post` metadata in `blog/<slug>.go`. Required: stable ASCII kebab-case `Slug`, `Title`, `Description`, `PublishedAt`, explicit `Status`, `Hero`, `Social`, and `Body`. Optional: `Excerpt` (otherwise Description), `ModifiedAt`, `Topics`, `Headings`.
3. Use actual RFC 3339 timestamps with local UTC offsets, e.g. `2026-10-24T09:00:00+02:00`. Do not guess publication timing or refresh modification dates for cosmetic edits.
4. Copy only intended attachments into `public/blog/<slug>/`. Supply optimized JPEG/PNG (or SVG for diagrams), dimensions, and descriptive alt text. Keep hero artwork **5:2**, ideally 2000×800; optional smaller `Hero.Variants` must match that ratio. Add a **JPEG/PNG share preview**, at least 1200×630, composed from the same artwork without stretching. A simple 1200×630 padded version works. The hero is decorative on the article; alt text is retained for metadata and meaningful images elsewhere.
5. For longer posts, define `[]Heading` once and reuse entries through `@sectionHeading(headings[0])`. IDs must be unique kebab-case. Put the same slice in `Post.Headings`; the table of contents stays in sync. Reading time and word count are automatic estimates.
6. Register the post once in `Posts()` in `blog/posts.go`. Start with `Status: Draft`; use `Published` only after review. No route, sitemap, RSS, or card edits are needed. Published posts become visible at their timestamp; drafts never do. All compiled content is in the binary—drafts are not a place for secrets.
7. Run `just test`, `go vet ./...`, and `just build`. Check `/blog` and the article at phone/desktop widths, keyboard navigation, links, images, and share metadata. Include regenerated `*_templ.go` files with the changes.

## Content and brand

- Start with the reader's question and a useful answer. Use a concise title, a specific description (roughly 140–160 characters when natural), and descriptive H2/H3 headings. No forced word count.
- Pick one primary search topic and use relevant phrasing naturally. **No meta-keywords or keyword-density targets**; Google ignores that tag. Optional `Topics` are short, reader-facing labels, not a ranking mechanism. They also connect related articles; no empty tag archives are generated.
- Use the existing Space Grotesk/Azeret Mono, warm paper, fine borders, and restrained blue. Do not add per-post CSS, giant color panels, tracking scripts, or a client-side chart library.
- Use `@Callout("Label")` and `@Figure(src, alt, caption, width, height)` as needed. The example is a component catalogue, **not a requirement to include everything** in every post.
- Tables need a caption, scoped headers, and `.article-table-scroll` with a labelled `role="region"` and `tabindex="0"`. Code belongs in `<pre><code>` with horizontal scrolling. Never use `templ.Raw` for imported prose or code.
- Figures need informative captions and alt text; decorative images use empty alt. Charts need a title, text explanation, units, provenance, and accessible data. SVG/HTML diagrams should show meaningful decisions and exceptions. No unlabelled fictional metrics or fabricated quotes, sources, or client results.
- Link to relevant site content and real sources. Resolve `[[wikilinks]]` to real URLs; unresolved notes stay flagged for review, not dead links. Convert Obsidian callouts to `Callout`, `![[attachments]]` to local figures, Mermaid to accessible static diagrams, and footnotes to working references/backlinks. Do not expose private vault paths.
- FAQs are optional editorial content, not a reason to add unsupported rich-result markup. Titles/descriptions/schema must accurately describe visible content. Keep the full title; never truncate it in code.

## Example lifecycle

`ExamplePost()` is commented out in `Posts()`: no public card or article route. Keep its Markdown, templ, metadata, and assets as an authoring reference; **do not register it for production**. Only re-enable it for a local preview when explicitly asked, then disable it before committing. `Status: Example` additionally enforces `noindex` and exclusion from sitemap/RSS during such previews.

## Agent prompt

> Convert `[path/to/post.md]` into a blog post. Follow `blog/README.md` and use `blog/example.templ` for supported content patterns. Preserve my text and resolve its attachments/links. Ask about missing material. Add metadata/assets and register it as a draft; do not publish, commit, push, or deploy without my approval.
