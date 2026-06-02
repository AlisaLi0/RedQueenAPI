# RedQueen NSFW Content Moderation API — Examples

Ready-to-run code examples for the **NSFW Content Moderation API** on RapidAPI.

> 🔗 **Get your API key:** [NSFW Content Moderation API on RapidAPI](https://rapidapi.com/bleujours/api/nsfw-content-moderation-api)

OpenAI-compatible content moderation for **text and images**. Classifies content
across 13 safety categories (sexual, hate, harassment, self-harm, violence,
illicit and more) and returns per-category boolean flags plus confidence scores.
A drop-in replacement for the OpenAI `/v1/moderations` endpoint — ideal for
filtering user-generated content on forums, social platforms and UGC apps.

## Quick start

1. Subscribe (free tier available) and grab your key from the
   [RapidAPI listing](https://rapidapi.com/bleujours/api/nsfw-content-moderation-api).
2. Set it as an environment variable:

   ```bash
   export RAPIDAPI_KEY="your-rapidapi-key"
   ```

3. Pick your language below and run the example.

| Language   | File                                         |
| ---------- | -------------------------------------------- |
| cURL       | [curl/moderate.sh](curl/moderate.sh)         |
| Python     | [python/moderate.py](python/moderate.py)     |
| JavaScript | [javascript/moderate.js](javascript/moderate.js) |
| PHP        | [php/moderate.php](php/moderate.php)         |
| Go         | [go/main.go](go/main.go)                     |

## Endpoints

| Method | Path               | Description                                   |
| ------ | ------------------ | --------------------------------------------- |
| `POST` | `/v1/moderations`  | Moderate one or more text and/or image inputs |
| `POST` | `/detect`          | Friendly alias of `/v1/moderations`           |
| `GET`  | `/v1/models`       | List the moderation model and capabilities    |
| `GET`  | `/health`          | Liveness + upstream status                     |

Base URL (via RapidAPI proxy): `https://nsfw-content-moderation-api.p.rapidapi.com`

Required headers on every request:

```
X-RapidAPI-Key:  <your key>
X-RapidAPI-Host: nsfw-content-moderation-api.p.rapidapi.com
Content-Type:    application/json
```

## Request shape

`input` accepts a string, an array of strings, or an array of typed content parts:

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
```

## Response shape

```jsonc
{
  "id": "modr-1a2b3c4d5e6f7a8b9c0d1e2f",
  "model": "redqueen-moderation-001",
  "results": [
    {
      "flagged": true,
      "type": "text",
      "categories": { "sexual": true, "violence": false, "hate": false /* … */ },
      "category_scores": { "sexual": 0.98, "violence": 0.01, "hate": 0.00 /* … */ }
    }
  ]
}
```

`flagged` is `true` when any category crosses the moderation threshold. Use the
per-category booleans in `categories` (and the raw `category_scores`) to apply
your own policy.

## License

MIT — see [LICENSE](LICENSE). Use these snippets freely in your own projects.
