// Moderate text and images with the NSFW Content Moderation API.
//
//   export RAPIDAPI_KEY="your-rapidapi-key"
//   node moderate.js
//
// Node 18+ (built-in fetch). Get your key:
// https://rapidapi.com/bleujours/api/nsfw-content-moderation-api

const HOST = "nsfw-content-moderation-api.p.rapidapi.com";
const BASE_URL = `https://${HOST}`;

const API_KEY = process.env.RAPIDAPI_KEY;
if (!API_KEY) {
  console.error("Set RAPIDAPI_KEY first (see the RapidAPI listing).");
  process.exit(1);
}

const headers = {
  "X-RapidAPI-Key": API_KEY,
  "X-RapidAPI-Host": HOST,
  "Content-Type": "application/json",
};

async function moderate(payload, path = "/v1/moderations") {
  const resp = await fetch(BASE_URL + path, {
    method: "POST",
    headers,
    body: JSON.stringify(payload),
  });
  if (!resp.ok) {
    throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
  }
  return resp.json();
}

async function main() {
  // 1) Moderate a single text string.
  const result = await moderate({ input: "explicit hardcore content all night" });
  const first = result.results[0];
  console.log("flagged:", first.flagged);
  const flagged = Object.entries(first.categories)
    .filter(([, hit]) => hit)
    .map(([name]) => name);
  console.log("categories:", flagged.join(", ") || "(none)");

  // 2) Batch: mix text and an image URL.
  const batch = await moderate({
    input: [
      { type: "text", text: "I love baking bread with my grandmother" },
      { type: "image_url", image_url: { url: "https://example.com/photo.jpg" } },
    ],
  });
  batch.results.forEach((item, i) => {
    console.log(`item ${i} (${item.type}): flagged=${item.flagged}`);
  });
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
