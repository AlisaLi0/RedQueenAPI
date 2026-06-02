<?php
// Moderate text with the NSFW Content Moderation API.
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

function moderate($baseUrl, $host, $apiKey, $payload, $path = "/v1/moderations") {
    $ch = curl_init($baseUrl . $path);
    curl_setopt_array($ch, [
        CURLOPT_POST => true,
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_HTTPHEADER => [
            "X-RapidAPI-Key: {$apiKey}",
            "X-RapidAPI-Host: {$host}",
            "Content-Type: application/json",
        ],
        CURLOPT_POSTFIELDS => json_encode($payload),
    ]);
    $body = curl_exec($ch);
    $status = curl_getinfo($ch, CURLINFO_HTTP_CODE);
    curl_close($ch);
    if ($status !== 200) {
        throw new Exception("HTTP {$status}: {$body}");
    }
    return json_decode($body, true);
}

// 1) Moderate a single text string.
$result = moderate($baseUrl, $host, $apiKey, ["input" => "explicit hardcore content all night"]);
$first = $result["results"][0];
echo "flagged: " . ($first["flagged"] ? "true" : "false") . "\n";
$flagged = array_keys(array_filter($first["categories"]));
echo "categories: " . (implode(", ", $flagged) ?: "(none)") . "\n";

// 2) Batch: mix text and an image URL.
$batch = moderate($baseUrl, $host, $apiKey, [
    "input" => [
        ["type" => "text", "text" => "I love baking bread with my grandmother"],
        ["type" => "image_url", "image_url" => ["url" => "https://example.com/photo.jpg"]],
    ],
]);
foreach ($batch["results"] as $i => $item) {
    echo "item {$i} ({$item['type']}): flagged=" . ($item["flagged"] ? "true" : "false") . "\n";
}
