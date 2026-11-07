# Cover images

Cover images (Eduza educations, social posts, etc.) are **rendered from HTML, not AI-generated**: every text, colour and position stays editable, and a re-render takes seconds. They use the blog's look: warm paper, a fine 48px grid, corner brackets, Space Grotesk/Azeret Mono and restrained blue/gold (tokens copied from `public/styles.css`).

Nothing here is used by the Go site or embedded in the binary.

```
tools/cover-image/
├── templates/                  one HTML file per image
│   └── eduza-upoznaj-programera.html
├── assets/
│   └── erik-cutout.png         transparent portrait, reused by templates
├── scripts/
│   ├── render.sh               HTML → PNG (headless Chrome)
│   ├── cutout.py               photo → transparent PNG (BiRefNet + edge decontamination)
│   └── check_edges.py          zoomed sheet for checking a cutout's edges
├── requirements.txt            Python deps for cutout/check scripts
└── out/                        rendered PNGs (git-ignored)
```

## Requirements

- macOS with Google Chrome (override the path with `CHROME=...`).
- Internet while rendering (fonts come from Google Fonts, as on the site).
- Python 3 only if you make a new cutout:

```bash
python3 -m venv tools/cover-image/.venv
tools/cover-image/.venv/bin/pip install -r tools/cover-image/requirements.txt
```

## Render an image

```bash
tools/cover-image/scripts/render.sh tools/cover-image/templates/eduza-upoznaj-programera.html
```

Writes `tools/cover-image/out/eduza-upoznaj-programera.png` at 1908×1136. Optional arguments: `[out.png] [width] [height] [scale]`. Use scale `2` for a 2× export. Width/height must match the `html, body` size in the template.

While editing, open the template in a browser at 100% zoom to preview it. The PNG from `render.sh` is the source of truth.

## Make a new image

1. Copy the closest template in `templates/` and give it a descriptive kebab-case name.
2. Set the canvas size in `html, body` (and the `svg.corners` viewBox/paths) to the target size.
3. Edit the copy. Typical layout: kicker → eyebrow → headline on the left, portrait card on the right.
4. Render, look at the PNG at full size **and** shrunk to ~500px wide (that is how cards and thumbnails show it), then adjust.

### Eduza

- Required: PNG/JPG, **1908×1136 px or larger, 1.68:1 landscape**.
- Eduza puts the image in a white frame and in small cards, so keep the headline large (≥ ~100px at 1908 wide) and keep the important content away from the edges.

### Template anatomy (`eduza-upoznaj-programera.html`)

| Piece | What to change |
| --- | --- |
| `.kicker` | Format label (blue pill) + meta line ("45 min · online · besplatno"). Wraps to two lines at 36px. |
| `.eyebrow` | Short blue mono line above the headline (education name). |
| `h1` | Main message. The blue `.dot` at the end is the brand accent. Shrink `font-size` for longer text. |
| `.byline` | Name / site, pinned to the bottom left. |
| `.frame` | Paper card with blue offset shadow. Runs off the bottom edge on purpose. |
| `.window` | Clips the portrait to the card's inner width, so the card sides cut the shoulders (hiding the photo's own crop edges) while the head breaks out above the card. If you move or resize `.frame`, update `.window`: `left` = frame left + 3px border, `width` = frame width − 6px. |
| `.person` | Portrait size (`width`) and horizontal position (`left`, relative to `.window`). Bottom-anchored to the canvas. |

## Make a new cutout

Use this when a new photo is needed. The current `assets/erik-cutout.png` comes from `~/Downloads/profil.JPG` (2388×2388).

```bash
tools/cover-image/.venv/bin/python -I tools/cover-image/scripts/cutout.py ~/Downloads/profil.JPG tools/cover-image/assets/erik-cutout.png
tools/cover-image/.venv/bin/python -I tools/cover-image/scripts/check_edges.py tools/cover-image/assets/erik-cutout.png tools/cover-image/out/edges.jpg
```

- `cutout.py` uses **BiRefNet-portrait** (via `rembg`) for the alpha mask. Its first run downloads the model (~1 GB to `~/.rembg/models`). It then runs **pymatting** foreground estimation to remove background colour that bleeds into the edges (e.g. green foliage in the hair). It saves at most 1600px on the longest side (`--size`), which is plenty for a 1908px-wide canvas. Use `--size 0` for full resolution, e.g. for 2× renders.
- `check_edges.py` writes a sheet of enlarged crops (hairline, hair, ear, neck/shoulder) on a dark and a light background. Look for halos, coloured fringe or jagged steps before using the cutout. Pass custom `x,y,size` boxes if the subject is framed differently.
- Compared alternatives: Apple Vision's subject lifting (`VNGenerateForegroundInstanceMaskRequest`, the "Lift subject" feature in Photos) was good but left a jagged step at the neck/shoulder on this photo. BiRefNet gave cleaner edges and softer hair. If a photo is hard (busy background in the hair), try `--model birefnet-general` or retake against a plain wall.
- Photos work best when the person faces the camera, is in sharp focus with even light, and the shoulders are visible. A crop at the photo's bottom/sides is fine: the template's card hides it.
