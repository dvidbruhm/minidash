#!/usr/bin/env bash
set -euo pipefail
mkdir -p internal/icons/collections web/static/vendor
for col in simple-icons lucide tabler; do
  curl -fsSL "https://cdn.jsdelivr.net/npm/@iconify-json/${col}/icons.json" -o "internal/icons/collections/${col}.json"
done
curl -fsSL https://unpkg.com/alpinejs@3.14.1/dist/cdn.min.js -o web/static/vendor/alpine.min.js
curl -fsSL https://cdn.jsdelivr.net/npm/sortablejs@1.15.3/Sortable.min.js -o web/static/vendor/sortable.min.js
echo done
