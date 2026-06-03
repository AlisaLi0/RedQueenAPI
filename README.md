# RedQueen Content Moderation — Examples

Ready-to-run code examples for RedQueen's two content-moderation APIs on RapidAPI.
Both share the same OpenAI-style surface and **one RapidAPI key** — subscribe to
whichever fits your use case (or both).

| # | API | Best for | Listing |
| - | --- | -------- | ------- |
| 1 | **NSFW Content Moderation API** | Fast, cheap, **image-only** safe/NSFW check | [nsfw-content-moderation-api](https://rapidapi.com/bleujours/api/nsfw-content-moderation-api) |
| 2 | **AI Content Moderation API** | Reasoning LLM, **text + images**, 13 categories | [ai-content-moderation-api](https://rapidapi.com/bleujours/api/ai-content-moderation-api) |

## Quick start

1. Subscribe (free tier available) to one or both APIs and grab your key.
2. Set it as an environment variable:

   ```bash
   export RAPIDAPI_KEY="your-rapidapi-key"
   ```

3. Pick your language below and run the example. Each example exercises **both**
   APIs end to end.

| Language   | File                                              |
| ---------- | ------------------------------------------------- |
| cURL       | [curl/moderate.sh](curl/moderate.sh)              |
| Python     | [python/moderate.py](python/moderate.py)          |
| JavaScript | [javascript/moderate.js](javascript/moderate.js)  |
| PHP        | [php/moderate.php](php/moderate.php)              |
| Go         | [go/main.go](go/main.go)                          |

---

## 1) NSFW Content Moderation API — fast, image-only

A lightweight vision classifier that returns a binary safe/NSFW verdict with a
confidence score in milliseconds. **Image input only** — sending text returns
HTTP `400`.

Base URL: `https://nsfw-content-moderation-api.p.rapidapi.com`

```
X-RapidAPI-Key:  <your key>
X-RapidAPI-Host: nsfw-content-moderation-api.p.rapidapi.com
Content-Type:    application/json
```

| Method | Path              | Description                                  |
| ------ | ----------------- | -------------------------------------------- |
| `POST` | `/v1/moderations` | Moderate one image (`image_url`/`image_b64`) |
| `POST` | `/detect`         | Friendly alias of `/v1/moderations`          |
| `GET`  | `/v1/models`      | Model id and capabilities                    |
| `GET`  | `/health`         | Liveness + upstream status                   |

**Request** — an image URL or base64-encoded bytes:

```jsonc
{ "image_url": "https://example.com/photo.jpg" }
// or
{ "image_b64": "iVBORw0KGgo..." }
```

**Response:**

```jsonc
{
  "id": "modr-ea4bc9e9370a45688120002b",
  "model": "redqueen-moderation-001-fast",
  "results": [
    {
      "flagged": false,
      "type": "image",
      "nsfw_score": 0.000195,
      "category_scores": { "normal": 0.999806, "nsfw": 0.000195 }
    }
  ],
  "latency_ms": 41
}
```

---

## 2) AI Content Moderation API — text + image, 13 categories

A reasoning LLM (Qwen3.6) that classifies **text and/or images** across 13 safety
categories and returns per-category booleans plus confidence scores. Drop-in
compatible with the OpenAI `/v1/moderations` shape.

Base URL: `https://ai-content-moderation-api.p.rapidapi.com`

```
X-RapidAPI-Key:  <your key>
X-RapidAPI-Host: ai-content-moderation-api.p.rapidapi.com
Content-Type:    application/json
```

| Method | Path              | Description                                   |
| ------ | ----------------- | --------------------------------------------- |
| `POST` | `/v1/moderations` | Moderate one or more text and/or image inputs |
| `POST` | `/detect`         | Friendly alias of `/v1/moderations`           |
| `GET`  | `/v1/models`      | Model id and capabilities                     |
| `GET`  | `/health`         | Liveness + upstream status                    |

**Request** — `input` accepts a string, an array of strings, or typed parts:

```jsonc
// 1) single string
{ "input": "some user text to check" }

// 2) array of strings (batch)
{ "input": ["first message", "second message"] }

// 3) typed parts — mix text and images
{
  "input": [
    { "type": "text", "text": "I love baking bread with my grandmother" },
    { "type": "image_url", "image_url": { "url": "https://example.com/photo.jpg" } }
  ]
}

// 4) image as a base64 data URL (e.g. a user upload)
{
  "input": [
    { "type": "image_url", "image_url": { "url": "data:image/png;base64,iVBORw0KGgo..." } }
  ]
}
```

**Response** — `categories` and `category_scores` always include all 13 keys:

```jsonc
{
  "id": "modr-f00de6539efc4afa83dbed1b",
  "model": "redqueen-moderation-001",
  "results": [
    {
      "flagged": true,
      "type": "text",
      "categories": { "harassment": true, "harassment/threatening": true, "violence": true /* … */ },
      "category_scores": { "harassment": 0.7, "harassment/threatening": 0.8, "violence": 0.6 /* … */ }
    }
  ]
}
```

The 13 categories: `sexual`, `sexual/minors`, `harassment`, `harassment/threatening`,
`hate`, `hate/threatening`, `illicit`, `illicit/violent`, `self-harm`,
`self-harm/intent`, `self-harm/instructions`, `violence`, `violence/graphic`.

`flagged` is `true` when any category crosses the moderation threshold. Use the
per-category booleans in `categories` (and the raw `category_scores`) to apply
your own policy.

---

## Rate limits

Each plan defines a monthly request quota. When you exceed it the RapidAPI proxy
returns **HTTP 429**. Check the response status and the optional `Retry-After`
header, then back off and retry — every example in this repo does this for you.

## License

MIT — see [LICENSE](LICENSE). Use these snippets freely in your own projects.
