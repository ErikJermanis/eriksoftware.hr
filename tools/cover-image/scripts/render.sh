#!/bin/zsh
# Render an HTML template to PNG with headless Google Chrome.
#
#   scripts/render.sh <template.html> [out.png] [width] [height] [scale]
#
# Defaults: out = out/<template-name>.png, 1908x1136 (Eduza), scale 1.
# Use scale 2 for a high-DPI export (e.g. 3816x2272).
set -euo pipefail

CHROME="${CHROME:-/Applications/Google Chrome.app/Contents/MacOS/Google Chrome}"
HERE="${0:A:h}"
ROOT="${HERE:h}"

[[ $# -ge 1 ]] || { echo "usage: $0 <template.html> [out.png] [width] [height] [scale]" >&2; exit 1; }

SRC="${1:A}"
OUT="${2:-$ROOT/out/${${1:t}:r}.png}"
OUT="${OUT:A}"
W="${3:-1908}"
H="${4:-1136}"
S="${5:-1}"

mkdir -p "${OUT:h}"
PROFILE="$(mktemp -d)"
trap 'rm -rf "$PROFILE"' EXIT

# --virtual-time-budget gives Google Fonts time to load before the screenshot.
"$CHROME" --headless=new --disable-gpu --hide-scrollbars \
  --force-device-scale-factor="$S" --window-size="$W,$H" \
  --virtual-time-budget=8000 --user-data-dir="$PROFILE" \
  --screenshot="$OUT" "file://$SRC" 2>/dev/null

echo "$OUT"
