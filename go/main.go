// Moderate text with the NSFW Content Moderation API.
//
//	export RAPIDAPI_KEY="your-rapidapi-key"
//	go run main.go
//
// Get your key: https://rapidapi.com/bleujours/api/nsfw-content-moderation-api
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const host = "nsfw-content-moderation-api.p.rapidapi.com"
const baseURL = "https://" + host

type moderationResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged    bool            `json:"flagged"`
		Type       string          `json:"type"`
		Categories map[string]bool `json:"categories"`
	} `json:"results"`
}

func moderate(apiKey string, payload any, path string) (*moderationResponse, error) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-RapidAPI-Key", apiKey)
	req.Header.Set("X-RapidAPI-Host", host)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, raw)
	}
	var out moderationResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func main() {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Set RAPIDAPI_KEY first (see the RapidAPI listing).")
		os.Exit(1)
	}

	// 1) Moderate a single text string.
	result, err := moderate(apiKey, map[string]any{
		"input": "explicit hardcore content all night",
	}, "/v1/moderations")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	first := result.Results[0]
	fmt.Println("flagged:", first.Flagged)
	for name, hit := range first.Categories {
		if hit {
			fmt.Println("category:", name)
		}
	}

	// 2) Batch: mix text and an image URL.
	batch, err := moderate(apiKey, map[string]any{
		"input": []any{
			map[string]any{"type": "text", "text": "I love baking bread with my grandmother"},
			map[string]any{"type": "image_url", "image_url": map[string]any{"url": "https://example.com/photo.jpg"}},
		},
	}, "/detect")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, item := range batch.Results {
		fmt.Printf("item %d (%s): flagged=%v\n", i, item.Type, item.Flagged)
	}
}
