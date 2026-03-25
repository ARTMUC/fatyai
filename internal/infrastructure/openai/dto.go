package openai

// chatRequest is the JSON body sent to the Chat Completions endpoint.
type chatRequest struct {
	Model          string        `json:"model"`
	Messages       []chatMessage `json:"messages"`
	MaxTokens      int           `json:"max_tokens"`
	ResponseFormat *respFormat   `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string        `json:"role"`
	Content interface{}   `json:"content"`
}

type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}

type imageURL struct {
	URL string `json:"url"`
}

type respFormat struct {
	Type string `json:"type"`
}

// chatResponse is the JSON body returned by the Chat Completions endpoint.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// AnalysisResponse is the structured result parsed from GPT-4o's JSON output.
type AnalysisResponse struct {
	Name         string  `json:"name"`
	CaloriesKcal float64 `json:"calories_kcal"`
	ProteinG     float64 `json:"protein_g"`
	FatG         float64 `json:"fat_g"`
	CarbsG       float64 `json:"carbs_g"`
	WeightG      float64 `json:"weight_g"`
	Estimated    bool    `json:"estimated"`
}
