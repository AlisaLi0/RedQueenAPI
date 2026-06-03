#!/usr/bin/env node
/*
 * RedQueen content moderation -- JavaScript (Node 18+) examples.
 *
 * Two complementary APIs share one RapidAPI key. Subscribe to whichever you need:
 *
 *   1) NSFW Content Moderation API  (fast, image-only safe/NSFW check)
 *      host: nsfw-content-moderation-api.p.rapidapi.com
 *      https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
 *
 *   2) AI Content Moderation API    (reasoning LLM, text + image, 13 categories)
 *      host: ai-content-moderation-api.p.rapidapi.com
 *      https://rapidapi.com/bleujours/api/ai-content-moderation-api
 *
 * Usage:  RAPIDAPI_KEY=your-key node moderate.js
 * Node 18+ has a built-in global fetch.
 */
const KEY = process.env.RAPIDAPI_KEY;
if (!KEY) {
  console.error("Set RAPIDAPI_KEY to your RapidAPI key");
  process.exit(1);
}

const NSFW_HOST = "nsfw-content-moderation-api.p.rapidapi.com";
const AI_HOST = "ai-content-moderation-api.p.rapidapi.com";

// 1x1 transparent PNG, used for the base64 image example.
const SAMPLE_PNG =
  "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==";

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

async function request(host, method, path, body, maxRetries = 3) {
  const url = `https://${host}${path}`;
  const headers = {
    "X-RapidAPI-Key": KEY,
    "X-RapidAPI-Host": host,
    "Content-Type": "application/json",
  };
  for (let attempt = 0; ; attempt++) {
    const resp = await fetch(url, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    if (resp.status === 429 && attempt < maxRetries) {
      const wait = Number(resp.headers.get("Retry-After")) || 2 ** attempt;
      await sleep(wait * 1000);
      continue;
    }
    if (!resp.ok) throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
    return resp.json();
  }
}

async function show(label, p) {
  console.log(`== ${label} ==`);
  console.log(JSON.stringify(await p, null, 2));
  console.log();
}

async function main() {
  console.log("### Product 1 -- NSFW Content Moderation API (fast, image-only)\n");
  await show("health", request(NSFW_HOST, "GET", "/health"));
  await show("models", request(NSFW_HOST, "GET", "/v1/models"));
  await show("image by URL", request(NSFW_HOST, "POST", "/v1/moderations", {
    image_url: "https://picsum.photos/id/237/300/300",
  }));
  await show("image by base64 (/detect)", request(NSFW_HOST, "POST", "/detect", {
    image_b64: SAMPLE_PNG,
  }));
  // This API is image-only. Sending {input: "text"} returns HTTP 400.

  console.log("### Product 2 -- AI Content Moderation API (text + image, 13 cat)\n");
  await show("health", request(AI_HOST, "GET", "/health"));
  await show("models", request(AI_HOST, "GET", "/v1/models"));
  await show("single text", request(AI_HOST, "POST", "/v1/moderations", {
    input: "I will hunt you down and hurt you",
  }));
  await show("batch strings", request(AI_HOST, "POST", "/v1/moderations", {
    input: ["hello there", "explicit hardcore content all night"],
  }));
  await show("text + image (/detect)", request(AI_HOST, "POST", "/detect", {
    input: [
      { type: "text", text: "check this" },
      { type: "image_url", image_url: { url: "https://picsum.photos/id/237/300/300" } },
    ],
  }));
  await show("image by base64 data URL", request(AI_HOST, "POST", "/v1/moderations", {
    input: [
      { type: "image_url", image_url: { url: `data:image/png;base64,${SAMPLE_PNG}` } },
    ],
  }));
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
