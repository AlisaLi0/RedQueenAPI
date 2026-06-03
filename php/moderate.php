<?php
// Moderate text and images with the NSFW Content Moderation API.
//
//   export RAPIDAPI_KEY="your-rapidapi-key"
//   php moderate.php
//
// Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api

$host = "nsfw-content-moderation-api.p.rapidapi.com";
$baseUrl = "https://{$host}";

$apiKey = getenv("RAPIDAPI_KEY");
if (!$apiKey) {
    fwrite(STDERR, "Set RAPIDAPI_KEY first (see the RapidAPI listing).\n");
    exit(1);
}

// A 1x1 PNG. Swap in your own bytes (e.g. file_get_contents("pic.jpg")).
$samplePng = base64_decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==");

// Build a base64 data URL the API accepts from raw bytes.
function dataUrl($bytes, $mime = "image/png") {
    return "data:{$mime};base64," . base64_encode($bytes);
}

// Call the API, retrying on HTTP 429 using the Retry-After header.
function apiRequest($method, $baseUrl, $host, $apiKey, $path, $payload = null, $maxRetries = 3) {
    for ($attempt = 0; ; $attempt++) {
        $ch = curl_init($baseUrl . $path);
        $opts = [
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_HEADER => true,
            CURLOPT_HTTPHEADER => [
                "X-RapidAPI-Key: {$apiKey}",
                "X-RapidAPI-Host: {$host}",
                "Content-Type: application/json",
            ],
        ];
        if ($payload !== null) {
            $opts[CURLOPT_POSTFIELDS] = json_encode($payload);
        }
        curl_setopt_array($ch, $opts);
        $raw = curl_exec($ch);
        $status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $headerSize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
        curl_close($ch);

        $rawHeaders = substr($raw, 0, $headerSize);
        $body = substr($raw, $headerSize);

        if ($status === 429 && $attempt < $maxRetries) {
            preg_match('/retry-after:\s*(\d+)/i', $rawHeaders, $m);
            $wait = isset($m[1]) ? (int)$m[1] : (1 << $attempt);
            fwrite(STDERR, "rate limited, retrying in {$wait}s...\n");
            sleep($wait);
            continue;
        }
        if ($status !== 200) {
            throw new Exception("HTTP {$status}: {$body}");
        }
        return json_decode($body, true);
    }
}

function moderate($baseUrl, $host, $apiKey, $payload, $path = "/v1/moderations") {
    return apiRequest("POST", $baseUrl, $host, $apiKey, $path, $payload);
}

// 0) Liveness + which model is serving.
echo "health: " . json_encode(apiRequest("GET", $baseUrl, $host, $apiKey, "/health")) . "\n";
echo "models: " . json_encode(apiRequest("GET", $baseUrl, $host, $apiKey, "/v1/models")) . "\n";

// 1) Moderate a single text string.
$result = moderate($baseUrl, $host, $apiKey, ["input" => "explicit hardcore content all night"]);
$first = $result["results"][0];
echo "flagged: " . ($first["flagged"] ? "true" : "false") . "\n";
$flagged = array_keys(array_filter($first["categories"]));
echo "categories: " . (implode(", ", $flagged) ?: "(none)") . "\n";

// 2) Batch of plain strings.
$strings = moderate($baseUrl, $host, $apiKey, ["input" => ["first message to check", "second message to check"]]);
foreach ($strings["results"] as $i => $item) {
    echo "string {$i}: flagged=" . ($item["flagged"] ? "true" : "false") . "\n";
}

// 3) Mix text and an image URL (via the /detect alias).
$batch = moderate($baseUrl, $host, $apiKey, [
    "input" => [
        ["type" => "text", "text" => "I love baking bread with my grandmother"],
        ["type" => "image_url", "image_url" => ["url" => "https://example.com/photo.jpg"]],
    ],
], "/detect");
foreach ($batch["results"] as $i => $item) {
    echo "item {$i} ({$item['type']}): flagged=" . ($item["flagged"] ? "true" : "false") . "\n";
}

// 4) Moderate a local image as a base64 data URL.
$img = moderate($baseUrl, $host, $apiKey, [
    "input" => [["type" => "image_url", "image_url" => ["url" => dataUrl($samplePng)]]],
]);
echo "image flagged: " . ($img["results"][0]["flagged"] ? "true" : "false") . "\n";
