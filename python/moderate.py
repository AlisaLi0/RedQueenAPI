#!/usr/bin/env python3
"""RedQueen content moderation -- Python examples.

Two complementary APIs share one RapidAPI key. Subscribe to whichever you need:

  1) NSFW Content Moderation API  (fast, image-only safe/NSFW check)
     host: nsfw-content-moderation-api.p.rapidapi.com
     https://rapidapi.com/bleujours/api/nsfw-content-moderation-api

  2) AI Content Moderation API    (reasoning LLM, text + image, 13 categories)
     host: ai-content-moderation-api.p.rapidapi.com
     https://rapidapi.com/bleujours/api/ai-content-moderation-api

Usage:  RAPIDAPI_KEY=your-key python moderate.py
Requires: requests  (pip install requests)
"""
import json
import os
import time

import requests

KEY = os.environ.get("RAPIDAPI_KEY")
if not KEY:
    raise SystemExit("Set RAPIDAPI_KEY to your RapidAPI key")

NSFW_HOST = "nsfw-content-moderation-api.p.rapidapi.com"
AI_HOST = "ai-content-moderation-api.p.rapidapi.com"

# 1x1 transparent PNG, used for the base64 image example.
SAMPLE_PNG = (
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGA"
    "hKmMIQAAAABJRU5ErkJggg=="
)


def request(host, method, path, body=None, max_retries=3):
    """Call the API, retrying on HTTP 429 using the Retry-After header."""
    url = f"https://{host}{path}"
    headers = {
        "X-RapidAPI-Key": KEY,
        "X-RapidAPI-Host": host,
        "Content-Type": "application/json",
    }
    for attempt in range(max_retries + 1):
        resp = requests.request(method, url, headers=headers, json=body, timeout=60)
        if resp.status_code == 429 and attempt < max_retries:
            wait = float(resp.headers.get("Retry-After", 2 ** attempt))
            time.sleep(wait)
            continue
        resp.raise_for_status()
        return resp.json()
    resp.raise_for_status()


def show(label, data):
    print(f"== {label} ==")
    print(json.dumps(data, indent=2, ensure_ascii=False))
    print()


def main():
    print("### Product 1 -- NSFW Content Moderation API (fast, image-only)\n")
    show("health", request(NSFW_HOST, "GET", "/health"))
    show("models", request(NSFW_HOST, "GET", "/v1/models"))
    show("image by URL", request(NSFW_HOST, "POST", "/v1/moderations",
         {"image_url": "https://picsum.photos/id/237/300/300"}))
    show("image by base64 (/detect)", request(NSFW_HOST, "POST", "/detect",
         {"image_b64": SAMPLE_PNG}))
    # This API is image-only. Sending {"input": "text"} returns HTTP 400.

    print("### Product 2 -- AI Content Moderation API (text + image, 13 cat)\n")
    show("health", request(AI_HOST, "GET", "/health"))
    show("models", request(AI_HOST, "GET", "/v1/models"))
    show("single text", request(AI_HOST, "POST", "/v1/moderations",
         {"input": "I will hunt you down and hurt you"}))
    show("batch strings", request(AI_HOST, "POST", "/v1/moderations",
         {"input": ["hello there", "explicit hardcore content all night"]}))
    show("text + image (/detect)", request(AI_HOST, "POST", "/detect", {"input": [
        {"type": "text", "text": "check this"},
        {"type": "image_url", "image_url": {"url": "https://picsum.photos/id/237/300/300"}},
    ]}))
    show("image by base64 data URL", request(AI_HOST, "POST", "/v1/moderations", {"input": [
        {"type": "image_url", "image_url": {"url": f"data:image/png;base64,{SAMPLE_PNG}"}},
    ]}))


if __name__ == "__main__":
    main()
