// Package openai provides an HTTP client for the OpenAI Chat Completions API
// with GPT-4o Vision support.
package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	mealdomain "github.com/artmuc/fatyai/internal/domain/meal"
)

const (
	defaultBaseURL = "https://api.openai.com/v1"
	systemPrompt   = `You are a nutritionist AI. Analyze meal photos and return ONLY valid JSON.
If a weighing scale is visible, use the exact weight shown.
Otherwise, estimate a typical portion weight.
Set "estimated" to true when you estimate the weight.`
	userPrompt = `Analyze this meal photo and return ONLY this JSON object (no markdown, no extra text):
{"name":"<dish name>","calories_kcal":<number>,"protein_g":<number>,"fat_g":<number>,"carbs_g":<number>,"weight_g":<number>,"estimated":<true|false>}`
)

// Client is an HTTP client for the OpenAI Chat Completions API.
// It implements the mealdomain.VisionAnalyzer port.
type Client struct {
	httpClient *http.Client
	apiKey     string
	model      string
	baseURL    string
}

// NewClient creates a new OpenAI client.
func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "gpt-4o"
	}
	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		model:      model,
		baseURL:    defaultBaseURL,
	}
}

// AnalyzeMeal implements mealdomain.VisionAnalyzer.
func (c *Client) AnalyzeMeal(ctx context.Context, imageBase64, mimeType string) (*mealdomain.VisionResult, error) {
	imageDataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, imageBase64)

	body := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{
				Role: "user",
				Content: []contentPart{
					{Type: "text", Text: userPrompt},
					{Type: "image_url", ImageURL: &imageURL{URL: imageDataURL}},
				},
			},
		},
		MaxTokens:      400,
		ResponseFormat: &respFormat{Type: "json_object"},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if chatResp.Error != nil {
		return nil, fmt.Errorf("openai error: %s", chatResp.Error.Message)
	}
	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	var analysis AnalysisResponse
	if err := json.Unmarshal([]byte(chatResp.Choices[0].Message.Content), &analysis); err != nil {
		return nil, fmt.Errorf("parse analysis json: %w", err)
	}

	return &mealdomain.VisionResult{
		Name:         analysis.Name,
		CaloriesKcal: analysis.CaloriesKcal,
		ProteinG:     analysis.ProteinG,
		FatG:         analysis.FatG,
		CarbsG:       analysis.CarbsG,
		WeightG:      analysis.WeightG,
		Estimated:    analysis.Estimated,
	}, nil
}
