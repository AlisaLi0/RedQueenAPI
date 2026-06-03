// Moderate text and images with the NSFW Content Moderation API.
//
//	export RAPIDAPI_KEY="your-rapidapi-key"
//	go run main.go
//
// Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const host = "nsfw-content-moderation-api.p.rapidapi.com"
const baseURL = "https://" + host

// A 1x1 PNG. Swap in your own bytes (e.g. os.ReadFile("pic.jpg")).
const samplePNGB64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="

type moderationResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged    bool            `json:"flagged"`
		Type       string          `json:"type"`
		Categories map[string]bool `json:"categories"`
	} `json:"results"`
}

var client = &http.Client{Timeout: 30 * time.Second}

// dataURL builds a base64 data URL the API accepts from raw image bytes.
func dataURL(b []byte, mime string) string {
	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(b)
}

// apiRequest calls the API, retrying on HTTP 429 using the Retry-After header.
func apiRequest(apiKey, method, path string, payload any) ([]byte, error) {
	const maxRetries = 3
	for attempt := 0; ; attempt++ {
		var reader io.Reader
		if payload != nil {
			body, _ := json.Marshal(payload)
			reader = bytes.NewReader(body)
		}
		req, err := http.NewRequest(method, baseURL+path, reader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-RapidAPI-Key", apiKey)
		req.Header.Set("X-RapidAPI-Host", host)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		raw, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && attempt < maxRetries {
			wait := 1 << attempt
			if ra, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil {
				wait = ra
			}
			fmt.Fprintf(os.Stderr, "rate limited, retrying in %ds...\n", wait)
			time.Sleep(time.Duration(wait) * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, raw)
		}
		return raw, nil
	}
}

func moderate(apiKey string, payload any, path string) (*moderationResponse, error) {
	raw, err := apiRequest(apiKey, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}
	var out moderationResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func must(label string, v *moderationResponse, err error) *moderationResponse {
	if err != nil {
		fmt.Fprintln(os.Stderr, label, err)
		os.Exit(1)
	}
	return v
}

func main() {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Set RAPIDAPI_KEY first (see the RapidAPI listing).")
		os.Exit(1)
	}

	// 0) Liveness + which model is serving.
	if raw, err := apiRequest(apiKey, http.MethodGet, "/health", nil); err == nil {
		fmt.Println("health:", string(raw))
	}
	if raw, err := apiRequest(apiKey, http.MethodGet, "/v1/models", nil); err == nil {
		fmt.Println("models:", string(raw))
	}

	// 1) Moderate a single text string.
	result := must("single:", moderate(apiKey, map[string]any{
		"input": "explicit hardcore content all night",
	}, "/v1/moderations"))
	first := result.Results[0]
	fmt.Println("flagged:", first.Flagged)
	for name, hit := range first.Categories {
		if hit {
			fmt.Println("category:", name)
		}
	}

	// 2) Batch of plain strings.
	strings := must("strings:", moderate(apiKey, map[string]any{
		"input": []any{"first message to check", "second message to check"},
	}, "/v1/moderations"))
	for i, item := range strings.Results {
		fmt.Printf("string %d: flagged=%v\n", i, item.Flagged)
	}

	// 3) Mix text and an image URL (via the /detect alias).
	batch := must("batch:", moderate(apiKey, map[string]any{
		"input": []any{
			map[string]any{"type": "text", "text": "I love baking bread with my grandmother"},
			map[string]any{"type": "image_url", "image_url": map[string]any{"url": "https://example.com/photo.jpg"}},
		},
	}, "/detect"))
	for i, item := range batch.Results {
		fmt.Printf("item %d (%s): flagged=%v\n", i, item.Type, item.Flagged)
	}

	// 4) Moderate a local image as a base64 data URL.
	png, _ := base64.StdEncoding.DecodeString(samplePNGB64)
	img := must("image:", moderate(apiKey, map[string]any{
		"input": []any{
			map[string]any{"type": "image_url", "image_url": map[string]any{"url": dataURL(png, "image/png")}},
		},
	}, "/v1/moderations"))
	fmt.Println("image flagged:", img.Results[0].Flagged)
}
