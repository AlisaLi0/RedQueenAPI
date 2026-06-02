#!/usr/bin/env python3
"""Moderate text and images with the NSFW Content Moderation API.

    pip install requests
    export RAPIDAPI_KEY="your-rapidapi-key"
    python moderate.py

Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
"""
import os
import sys

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


def moderate(payload, path="/v1/moderations"):
    resp = requests.post(BASE_URL + path, headers=HEADERS, json=payload, timeout=30)
    resp.raise_for_status()
    return resp.json()


def main():
    # 1) Moderate a single text string.
    result = moderate({"input": "explicit hardcore content all night"})
    first = result["results"][0]
    print("flagged:", first["flagged"])
    flagged_cats = [name for name, hit in first["categories"].items() if hit]
    print("categories:", ", ".join(flagged_cats) or "(none)")

    # 2) Batch: mix text and an image URL.
    batch = moderate({
        "input": [
            {"type": "text", "text": "I love baking bread with my grandmother"},
            {"type": "image_url", "image_url": {"url": "https://example.com/photo.jpg"}},
        ]
    })
    for i, item in enumerate(batch["results"]):
        print(f"item {i} ({item['type']}): flagged={item['flagged']}")


if __name__ == "__main__":
    main()
