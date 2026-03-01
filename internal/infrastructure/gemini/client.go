package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"

	"steam-bubbles/internal/domain"
)

const promptTemplate = `I enjoy these Steam games: %s.

Based on these games, suggest exactly 8 other games I might love.
Provide the result as a JSON array of objects.
Each object must have exactly three string fields:
- "name": The game name
- "reason": A short sentence explaining why it fits my taste
- "platform": Platform or where to get it (e.g. Steam, Switch, PlayStation)

Return ONLY the raw JSON array starting with [ and ending with ]. Do NOT wrap it in markdown block quotes.`

type Client struct {
	apiKey string
	model  string
}

func NewClient(apiKey, model string) *Client {
	return &Client{apiKey: apiKey, model: model}
}

func (c *Client) Recommend(ctx context.Context, games []domain.Game) ([]domain.Recommendation, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  c.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("gemini: client creation failed: %w", err)
	}

	names := make([]string, len(games))
	for i, g := range games {
		names[i] = g.Name
	}

	prompt := fmt.Sprintf(promptTemplate, strings.Join(names, ", "))

	result, err := client.Models.GenerateContent(ctx, c.model, genai.Text(prompt), &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
	})
	if err != nil {
		return nil, fmt.Errorf("gemini: content generation failed: %w", err)
	}

	jsonText := string(result.Text())

	// Strip out markdown code block if model still outputs it despite instructions
	jsonText = strings.TrimSpace(jsonText)
	if strings.HasPrefix(jsonText, "```") {
		lines := strings.Split(jsonText, "\n")
		if len(lines) > 2 {
			jsonText = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var recs []domain.Recommendation
	if err := json.Unmarshal([]byte(jsonText), &recs); err != nil {
		return nil, fmt.Errorf("gemini: failed to parse JSON array: %w\nRaw text: %s", err, jsonText)
	}

	return recs, nil
}
