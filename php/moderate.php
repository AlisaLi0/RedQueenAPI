<?php
/*
 * RedQueen content moderation -- PHP examples.
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
 * Usage:  RAPIDAPI_KEY=your-key php moderate.php
 * Requires the cURL extension (ext-curl).
 */

$KEY = getenv('RAPIDAPI_KEY');
if (!$KEY) {
    fwrite(STDERR, "Set RAPIDAPI_KEY to your RapidAPI key\n");
    exit(1);
}

const NSFW_HOST = 'nsfw-content-moderation-api.p.rapidapi.com';
const AI_HOST   = 'ai-content-moderation-api.p.rapidapi.com';

// 1x1 transparent PNG, used for the base64 image example.
const SAMPLE_PNG = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==';

function api_request(string $host, string $method, string $path, ?array $body = null, int $maxRetries = 3) {
    global $KEY;
    $url = "https://{$host}{$path}";
    for ($attempt = 0; ; $attempt++) {
        $ch = curl_init($url);
        $headers = [
            "X-RapidAPI-Key: {$KEY}",
            "X-RapidAPI-Host: {$host}",
            'Content-Type: application/json',
        ];
        curl_setopt_array($ch, [
            CURLOPT_CUSTOMREQUEST  => $method,
            CURLOPT_HTTPHEADER     => $headers,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_HEADER         => true,
            CURLOPT_TIMEOUT        => 60,
        ]);
        if ($body !== null) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($body));
        }
        $raw = curl_exec($ch);
        $status = curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
        $hsize = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
        curl_close($ch);

        $rawHeaders = substr($raw, 0, $hsize);
        $bodyText = substr($raw, $hsize);

        if ($status === 429 && $attempt < $maxRetries) {
            $wait = 2 ** $attempt;
            if (preg_match('/retry-after:\s*(\d+)/i', $rawHeaders, $m)) {
                $wait = (int) $m[1];
            }
            sleep($wait);
            continue;
        }
        if ($status >= 400) {
            throw new RuntimeException("HTTP {$status}: {$bodyText}");
        }
        return json_decode($bodyText, true);
    }
}

function show(string $label, $data): void {
    echo "== {$label} ==\n";
    echo json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES | JSON_UNESCAPED_UNICODE) . "\n\n";
}

echo "### Product 1 -- NSFW Content Moderation API (fast, image-only)\n\n";
show('health', api_request(NSFW_HOST, 'GET', '/health'));
show('models', api_request(NSFW_HOST, 'GET', '/v1/models'));
show('image by URL', api_request(NSFW_HOST, 'POST', '/v1/moderations',
    ['image_url' => 'https://picsum.photos/id/237/300/300']));
show('image by base64 (/detect)', api_request(NSFW_HOST, 'POST', '/detect',
    ['image_b64' => SAMPLE_PNG]));
// This API is image-only. Sending ['input' => 'text'] returns HTTP 400.

echo "### Product 2 -- AI Content Moderation API (text + image, 13 cat)\n\n";
show('health', api_request(AI_HOST, 'GET', '/health'));
show('models', api_request(AI_HOST, 'GET', '/v1/models'));
show('single text', api_request(AI_HOST, 'POST', '/v1/moderations',
    ['input' => 'I will hunt you down and hurt you']));
show('batch strings', api_request(AI_HOST, 'POST', '/v1/moderations',
    ['input' => ['hello there', 'explicit hardcore content all night']]));
show('text + image (/detect)', api_request(AI_HOST, 'POST', '/detect', ['input' => [
    ['type' => 'text', 'text' => 'check this'],
    ['type' => 'image_url', 'image_url' => ['url' => 'https://picsum.photos/id/237/300/300']],
]]));
show('image by base64 data URL', api_request(AI_HOST, 'POST', '/v1/moderations', ['input' => [
    ['type' => 'image_url', 'image_url' => ['url' => 'data:image/png;base64,' . SAMPLE_PNG]],
]]));
