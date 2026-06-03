#!/usr/bin/env python3
"""Moderate text and images with the NSFW Content Moderation API.

    pip install requests
    export RAPIDAPI_KEY="your-rapidapi-key"
    python moderate.py

Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
"""
import base64
import os
import sys
import time

import requests

HOST = "nsfw-content-moderation-api.p.rapidapi.com"
BASE_URL = f"https://{HOST}"

API_KEY = os.environ.get("RAPIDAPI_KEY")
if not API_KEY:
    sys.exit("Set RAPIDAPI_KEY first (see https://rapidapi.com/bleujours/api/nsfw-content-moderation-api)")

HEADERS = {
    "X-RapidAPI-Key": API_KEY,
    "X-RapidAPI-Host": HOST,
    "Content-Type": "application/json",
}

# A 1x1 PNG. Swap in your own bytes (e.g. open("pic.jpg", "rb").read()).
SAMPLE_PNG = base64.b64decode(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
)


def data_url(image_bytes, mime="image/png"):
    """Turn raw image bytes into a base64 data URL the API accepts."""
    b64 = base64.b64encode(image_bytes).decode()
    return f"data:{mime};base64,{b64}"


def request(method, path, payload=None, max_retries=3):
    """Call the API, retrying on HTTP 429 using the Retry-After header."""
    for attempt in range(max_retries + 1):
        resp = requests.request(method, BASE_URL + path, headers=HEADERS, json=payload, timeout=30)
        if resp.status_code == 429 and attempt < max_retries:
            wait = float(resp.headers.get("Retry-After", 2 ** attempt))
            print(f"rate limited, retrying in {wait}s...", file=sys.stderr)
            time.sleep(wait)
            continue
        resp.raise_for_status()
        return resp.json()


def moderate(payload, path="/v1/moderations"):
    return request("POST", path, payload)


def main():
    # 0) Liveness + which model is serving.
    print("health:", request("GET", "/health"))
    print("models:", request("GET", "/v1/models"))

    # 1) Moderate a single text string.
    result = moderate({"input": "explicit hardcore content all night"})
    first = result["results"][0]
    print("flagged:", first["flagged"])
    flagged_cats = [name for name, hit in first["categories"].items() if hit]
    print("categories:", ", ".join(flagged_cats) or "(none)")

    # 2) Batch of plain strings.
    strings = moderate({"input": ["first message to check", "second message to check"]})
    for i, item in enumerate(strings["results"]):
        print(f"string {i}: flagged={item['flagged']}")

    # 3) Mix text and an image URL (via the /detect alias).
    batch = moderate(
        {
            "input": [
                {"type": "text", "text": "I love baking bread with my grandmother"},
                {"type": "image_url", "image_url": {"url": "https://example.com/photo.jpg"}},
            ]
        },
        path="/detect",
    )
    for i, item in enumerate(batch["results"]):
        print(f"item {i} ({item['type']}): flagged={item['flagged']}")

    # 4) Moderate a local image as a base64 data URL.
    img = moderate({"input": [{"type": "image_url", "image_url": {"url": data_url(SAMPLE_PNG)}}]})
    print("image flagged:", img["results"][0]["flagged"])


if __name__ == "__main__":
    main()
