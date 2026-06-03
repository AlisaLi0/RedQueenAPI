#!/usr/bin/env bash
#
# RedQueen content moderation — curl examples
#
# Two complementary APIs share one RapidAPI key. Subscribe to whichever you need:
#
#   1) NSFW Content Moderation API  (fast, image-only safe/NSFW check)
#      host: nsfw-content-moderation-api.p.rapidapi.com
#      https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
#
#   2) AI Content Moderation API    (reasoning LLM, text + image, 13 categories)
#      host: ai-content-moderation-api.p.rapidapi.com
#      https://rapidapi.com/bleujours/api/ai-content-moderation-api
#
# Usage:  RAPIDAPI_KEY=your-key ./moderate.sh
set -euo pipefail

KEY="${RAPIDAPI_KEY:?Set RAPIDAPI_KEY to your RapidAPI key}"
NSFW_HOST="nsfw-content-moderation-api.p.rapidapi.com"
AI_HOST="ai-content-moderation-api.p.rapidapi.com"

# 1x1 transparent PNG, used for the base64 image example.
SAMPLE_PNG="iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="

nsfw() { curl -sS "https://${NSFW_HOST}$1" -H "X-RapidAPI-Key: ${KEY}" -H "X-RapidAPI-Host: ${NSFW_HOST}" "${@:2}"; echo; }
ai()   { curl -sS "https://${AI_HOST}$1"   -H "X-RapidAPI-Key: ${KEY}" -H "X-RapidAPI-Host: ${AI_HOST}"   "${@:2}"; echo; }

echo "############################################################"
echo "# Product 1 — NSFW Content Moderation API (fast, image-only)"
echo "############################################################"

echo "== GET /health =="
nsfw /health

echo "== GET /v1/models =="
nsfw /v1/models

echo "== POST /v1/moderations  (image by URL) =="
nsfw /v1/moderations -X POST -H 'Content-Type: application/json' \
  -d '{"image_url":"https://picsum.photos/id/237/300/300"}'

echo "== POST /detect  (image by base64) =="
nsfw /detect -X POST -H 'Content-Type: application/json' \
  -d "{\"image_b64\":\"${SAMPLE_PNG}\"}"
# Returns: {"results":[{"flagged":false,"type":"image","nsfw_score":0.0002,
#                       "category_scores":{"normal":0.9998,"nsfw":0.0002}}]}
# NOTE: this API is image-only. Sending {"input":"some text"} returns HTTP 400.

echo
echo "############################################################"
echo "# Product 2 — AI Content Moderation API (text + image, 13 cat)"
echo "############################################################"

echo "== GET /health =="
ai /health

echo "== GET /v1/models =="
ai /v1/models

echo "== POST /v1/moderations  (single text) =="
ai /v1/moderations -X POST -H 'Content-Type: application/json' \
  -d '{"input":"I will hunt you down and hurt you"}'

echo "== POST /v1/moderations  (batch: array of strings) =="
ai /v1/moderations -X POST -H 'Content-Type: application/json' \
  -d '{"input":["hello there","explicit hardcore content all night"]}'

echo "== POST /detect  (text + image in one call) =="
ai /detect -X POST -H 'Content-Type: application/json' \
  -d '{"input":[{"type":"text","text":"check this"},{"type":"image_url","image_url":{"url":"https://picsum.photos/id/237/300/300"}}]}'

echo "== POST /v1/moderations  (image by base64 data URL) =="
ai /v1/moderations -X POST -H 'Content-Type: application/json' \
  -d "{\"input\":[{\"type\":\"image_url\",\"image_url\":{\"url\":\"data:image/png;base64,${SAMPLE_PNG}\"}}]}"
# Returns 13-category categories{} + category_scores{} per result.

# Rate limits: on HTTP 429, back off and retry using the Retry-After header.
