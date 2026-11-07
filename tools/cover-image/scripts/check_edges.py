"""Zoomed edge check for a cutout: composites it on a dark and a light background and
saves 4 enlarged crops per row, so halos, green fringe or jagged edges are easy to spot.

    python scripts/check_edges.py <cutout.png> <out.jpg> [x,y,size ...]

Crop boxes are in cutout pixels (top-left x, y and square size). The defaults suit a
head-and-shoulders portrait: left hairline, top-right hair, right ear, neck/shoulder.
"""
import sys

from PIL import Image

BACKGROUNDS = [(42, 30, 90), (240, 237, 228)]
TILE = 600


def main():
    cut = Image.open(sys.argv[1]).convert("RGBA")
    out = sys.argv[2]
    w = cut.width
    if len(sys.argv) > 3:
        boxes = [tuple(int(v) for v in s.split(",")) for s in sys.argv[3:]]
    else:  # fractions of width, matching a centred head-and-shoulders shot
        boxes = [(int(w * f[0]), int(w * f[1]), int(w * f[2])) for f in
                 [(.29, .10, .13), (.52, .10, .13), (.63, .34, .13), (.21, .71, .17)]]

    sheet = Image.new("RGB", (TILE * len(boxes), TILE * len(BACKGROUNDS)))
    for r, bg in enumerate(BACKGROUNDS):
        comp = Image.new("RGBA", cut.size, bg + (255,))
        comp.alpha_composite(cut)
        for c, (x, y, s) in enumerate(boxes):
            tile = comp.crop((x, y, x + s, y + s)).resize((TILE, TILE), Image.NEAREST)
            sheet.paste(tile.convert("RGB"), (c * TILE, r * TILE))
    sheet.save(out, quality=88)
    print(out)


if __name__ == "__main__":
    main()
