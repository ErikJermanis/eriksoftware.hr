"""Cut a person out of a photo and save a transparent PNG.

    python scripts/cutout.py <photo> <out.png> [--model birefnet-portrait] [--size 1600]

1. BiRefNet-portrait (via rembg) predicts a soft alpha mask. It keeps fine hair and
   gives cleaner neck/shoulder edges than Apple's Vision subject lifting.
2. pymatting's estimate_foreground_ml "decontaminates" edge pixels, removing the
   original background colour (e.g. green foliage) that bleeds into hair edges.
3. The result is optionally downscaled (--size, longest side) to keep the repo small.

First run downloads the model (~1 GB) to ~/.rembg/models.
"""
import argparse

import numpy as np
from PIL import Image, ImageOps
from pymatting import estimate_foreground_ml
from rembg import new_session, remove


def main():
    p = argparse.ArgumentParser()
    p.add_argument("photo")
    p.add_argument("out")
    p.add_argument("--model", default="birefnet-portrait")
    p.add_argument("--size", type=int, default=1600, help="longest side of output; 0 = keep")
    p.add_argument("--haze", type=float, default=0.03, help="alpha below this becomes fully transparent")
    a = p.parse_args()

    img = ImageOps.exif_transpose(Image.open(a.photo)).convert("RGB")
    mask = remove(img, session=new_session(a.model), only_mask=True).convert("L")
    mask = mask.resize(img.size, Image.LANCZOS)

    alpha = np.asarray(mask, dtype=np.float64) / 255.0
    alpha = np.clip((alpha - a.haze) / (1 - 2 * a.haze), 0, 1)  # drop faint haze, keep soft edge
    rgb = np.asarray(img, dtype=np.float64) / 255.0
    fg = estimate_foreground_ml(rgb, alpha)

    rgba = np.dstack([np.clip(fg, 0, 1), alpha])
    out = Image.fromarray((rgba * 255 + 0.5).astype(np.uint8), "RGBA")
    if a.size and max(out.size) > a.size:
        out.thumbnail((a.size, a.size), Image.LANCZOS)
    out.save(a.out, optimize=True)
    print(a.out, out.size)


if __name__ == "__main__":
    main()
