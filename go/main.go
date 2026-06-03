// RedQueen content moderation -- Go examples.
//
// Two complementary APIs share one RapidAPI key. Subscribe to whichever you need:
//
//	1) NSFW Content Moderation API  (fast, image-only safe/NSFW check)
//	   host: nsfw-content-moderation-api.p.rapidapi.com
//	   https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
//
//	2) AI Content Moderation API    (reasoning LLM, text + image, 13 categories)
//	   host: ai-content-moderation-api.p.rapidapi.com
//	   https://rapidapi.com/bleujours/api/ai-content-moderation-api
//
// Usage:  RAPIDAPI_KEY=your-key go run main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	nsfwHost = "nsfw-content-moderation-api.p.rapidapi.com"
	aiHost   = "ai-content-moderation-api.p.rapidapi.com"
	// samplePNG is a 1x1 transparent PNG used for the base64 image example.
	samplePNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
)

var key = os.Getenv("RAPIDAPI_KEY")

func apiRequest(host, method, path string, body any) (map[string]any, error) {
	url := "https://" + host + path
	maxRetries := 3
	for attempt := 0; ; attempt++ {
		var reader io.Reader
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reader = bytes.NewReader(b)
		}
		req, err := http.NewRequest(method, url, reader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-RapidAPI-Key", key)
		req.Header.Set("X-RapidAPI-Host", host)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		data, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			wait := 1 << attempt
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if n, err := strconv.Atoi(ra); err == nil {
					wait = n
				}
			}
			time.Sleep(time.Duration(wait) * time.Second)
			continue
		}
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
		}
		var out map[string]any
		if err := json.Unmarshal(data, &out); err != nil {
			return nil, err
		}
		return out, nil
	}
}

func show(label string, data map[string]any, err error) {
	if err != nil {
		fmt.Printf("== %s == ERROR: %v\n\n", label, err)
		return
	}
	pretty, _ := json.MarshalIndent(data, "", "  ")
	fmt.Printf("== %s ==\n%s\n\n", label, pretty)
}

func main() {
	if key == "" {
		fmt.Fprintln(os.Stderr, "Set RAPIDAPI_KEY to your RapidAPI key")
		os.Exit(1)
	}

	fmt.Println("### Product 1 -- NSFW Content Moderation API (fast, image-only)")
	fmt.Println()
	show(show3(apiRequest(nsfwHost, "GET", "/health", nil), "health"))
	show(show3(apiRequest(nsfwHost, "GET", "/v1/models", nil), "models"))
	show(show3(apiRequest(nsfwHost, "POST", "/v1/moderations",
		map[string]any{"image_url": "https://picsum.photos/id/237/300/300"}), "image by URL"))
	show(show3(apiRequest(nsfwHost, "POST", "/detect",
		map[string]any{"image_b64": samplePNG}), "image by base64 (/detect)"))
	// This API is image-only. Sending {"input":"text"} returns HTTP 400.

	fmt.Println("### Product 2 -- AI Content Moderation API (text + image, 13 cat)")
	fmt.Println()
	show(show3(apiRequest(aiHost, "GET", "/health", nil), "health"))
	show(show3(apiRequest(aiHost, "GET", "/v1/models", nil), "models"))
	show(show3(apiRequest(aiHost, "POST", "/v1/moderations",
		map[string]any{"input": "I will hunt you down and hurt you"}), "single text"))
	show(show3(apiRequest(aiHost, "POST", "/v1/moderations",
		map[string]any{"input": []string{"hello there", "explicit hardcore content all night"}}), "batch strings"))
	show(show3(apiRequest(aiHost, "POST", "/detect", map[string]any{"input": []any{
		map[string]any{"type": "text", "text": "check this"},
		map[string]any{"type": "image_url", "image_url": map[string]any{"url": "https://picsum.photos/id/237/300/300"}},
	}}), "text + image (/detect)"))
	show(show3(apiRequest(aiHost, "POST", "/v1/moderations", map[string]any{"input": []any{
		map[string]any{"type": "image_url", "image_url": map[string]any{"url": "data:image/png;base64," + samplePNG}},
	}}), "image by base64 data URL"))
}

// show3 adapts (data, err) plus a label into the argument order show expects.
func show3(data map[string]any, err error, label string) (string, map[string]any, error) {
	return label, data, err
}
