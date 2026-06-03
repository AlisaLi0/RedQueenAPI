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

// A 1x1 PNG. Swap in your own bytes (e.g. fs.readFileSync("pic.jpg")).
const SAMPLE_PNG_B64 =
  "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==";

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

// Build a base64 data URL the API accepts from raw bytes.
function dataUrl(base64, mime = "image/png") {
  return `data:${mime};base64,${base64}`;
}

// Call the API, retrying on HTTP 429 using the Retry-After header.
async function request(method, path, payload, maxRetries = 3) {
  for (let attempt = 0; ; attempt++) {
    const resp = await fetch(BASE_URL + path, {
      method,
      headers,
      body: payload ? JSON.stringify(payload) : undefined,
    });
    if (resp.status === 429 && attempt < maxRetries) {
      const wait = Number(resp.headers.get("Retry-After")) || 2 ** attempt;
      console.error(`rate limited, retrying in ${wait}s...`);
      await sleep(wait * 1000);
      continue;
    }
    if (!resp.ok) {
      throw new Error(`HTTP ${resp.status}: ${await resp.text()}`);
    }
    return resp.json();
  }
}

const moderate = (payload, path = "/v1/moderations") => request("POST", path, payload);

async function main() {
  // 0) Liveness + which model is serving.
  console.log("health:", await request("GET", "/health"));
  console.log("models:", await request("GET", "/v1/models"));

  // 1) Moderate a single text string.
  const result = await moderate({ input: "explicit hardcore content all night" });
  const first = result.results[0];
  console.log("flagged:", first.flagged);
  const flagged = Object.entries(first.categories)
    .filter(([, hit]) => hit)
    .map(([name]) => name);
  console.log("categories:", flagged.join(", ") || "(none)");

  // 2) Batch of plain strings.
  const strings = await moderate({ input: ["first message to check", "second message to check"] });
  strings.results.forEach((item, i) => console.log(`string ${i}: flagged=${item.flagged}`));

  // 3) Mix text and an image URL (via the /detect alias).
  const batch = await moderate(
    {
      input: [
        { type: "text", text: "I love baking bread with my grandmother" },
        { type: "image_url", image_url: { url: "https://example.com/photo.jpg" } },
      ],
    },
    "/detect"
  );
  batch.results.forEach((item, i) => {
    console.log(`item ${i} (${item.type}): flagged=${item.flagged}`);
  });

  // 4) Moderate a local image as a base64 data URL.
  const img = await moderate({
    input: [{ type: "image_url", image_url: { url: dataUrl(SAMPLE_PNG_B64) } }],
  });
  console.log("image flagged:", img.results[0].flagged);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
