#!/usr/bin/env bash
# Moderate text with the NSFW Content Moderation API.
#
#   export RAPIDAPI_KEY="your-rapidapi-key"
#   ./moderate.sh
#
# Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
set -euo pipefail

: "${RAPIDAPI_KEY:?Set RAPIDAPI_KEY first}"
HOST="nsfw-content-moderation-api.p.rapidapi.com"

# --- 1) Moderate a single text string -------------------------------------
curl -s -X POST "https://${HOST}/v1/moderations" \
  -H "X-RapidAPI-Key: ${RAPIDAPI_KEY}" \
  -H "X-RapidAPI-Host: ${HOST}" \
  -H "Content-Type: application/json" \
  -d '{"input": "explicit hardcore content all night"}'
echo

# --- 2) Batch text + image via the /detect alias --------------------------
curl -s -X POST "https://${HOST}/detect" \
  -H "X-RapidAPI-Key: ${RAPIDAPI_KEY}" \
  -H "X-RapidAPI-Host: ${HOST}" \
  -H "Content-Type: application/json" \
  -d '{
        "input": [
          { "type": "text", "text": "I love baking bread with my grandmother" },
          { "type": "image_url", "image_url": { "url": "https://example.com/photo.jpg" } }
        ]
      }'
echo
