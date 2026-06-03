#!/usr/bin/env bash
# Moderate text and images with the NSFW Content Moderation API.
#
#   export RAPIDAPI_KEY="your-rapidapi-key"
#   ./moderate.sh
#
# Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
set -euo pipefail

: "${RAPIDAPI_KEY:?Set RAPIDAPI_KEY first}"
HOST="nsfw-content-moderation-api.p.rapidapi.com"

auth=(-H "X-RapidAPI-Key: ${RAPIDAPI_KEY}" -H "X-RapidAPI-Host: ${HOST}")
json=(-H "Content-Type: application/json")

# --- 1) Health check ------------------------------------------------------
echo "# health"
curl -s "${auth[@]}" "https://${HOST}/health"
echo

# --- 2) List the moderation model -----------------------------------------
echo "# models"
curl -s "${auth[@]}" "https://${HOST}/v1/models"
echo

# --- 3) Moderate a single text string -------------------------------------
echo "# single text"
curl -s -X POST "https://${HOST}/v1/moderations" "${auth[@]}" "${json[@]}" \
  -d '{"input": "explicit hardcore content all night"}'
echo

# --- 4) Batch of plain strings --------------------------------------------
echo "# batch strings"
curl -s -X POST "https://${HOST}/v1/moderations" "${auth[@]}" "${json[@]}" \
  -d '{"input": ["first message to check", "second message to check"]}'
echo

# --- 5) Mixed text + image URL via the /detect alias ----------------------
echo "# text + image url"
curl -s -X POST "https://${HOST}/detect" "${auth[@]}" "${json[@]}" \
  -d '{
        "input": [
          { "type": "text", "text": "I love baking bread with my grandmother" },
          { "type": "image_url", "image_url": { "url": "https://example.com/photo.jpg" } }
        ]
      }'
echo

# --- 6) Image as a base64 data URL (e.g. a user upload) -------------------
# data:<mime>;base64,<...>  — swap in your own image bytes.
echo "# base64 image"
IMG_DATA_URL="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
curl -s -X POST "https://${HOST}/v1/moderations" "${auth[@]}" "${json[@]}" \
  -d "{\"input\": [{\"type\": \"image_url\", \"image_url\": {\"url\": \"${IMG_DATA_URL}\"}}]}"
echo

# Note on rate limits: when you exceed your plan quota the proxy returns
# HTTP 429. Inspect the response status and the Retry-After header, then back off.
