package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/razorpay/razorpay-cli/config"
	"github.com/spf13/cobra"
)

const (
	AnthropicAPIURL    = "https://api.anthropic.com/v1/messages"
	AnthropicVersion   = "2023-06-01"
	DefaultClaudeModel = "claude-haiku-4-5-20251001"

	GeminiAPIURL       = "https://generativelanguage.googleapis.com/v1beta/models/gemini-3.6-flash:generateContent"
	DefaultGeminiModel = "gemini-3.6-flash"
)

type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	rootCmd    *cobra.Command
}

// Anthropic Request/Response types
type anthropicMessageRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicResponse struct {
	ID      string                  `json:"id"`
	Type    string                  `json:"type"`
	Role    string                  `json:"role"`
	Content []anthropicContentBlock `json:"content"`
	Error   *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Gemini Request/Response types
type geminiRequest struct {
	SystemInstruction *geminiContent         `json:"system_instruction,omitempty"`
	Contents          []geminiContent        `json:"contents"`
	GenerationConfig  map[string]interface{} `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func NewClient() *Client {
	return NewClientWithKey(config.AIApiKey())
}

func NewClientWithKey(apiKey string) *Client {
	return &Client{
		apiKey:     strings.TrimSpace(apiKey),
		model:      DefaultClaudeModel,
		httpClient: &http.Client{Timeout: 20 * time.Second},
		rootCmd:    GetGlobalRootCommand(),
	}
}

// SetRootCommand sets the root command used for dynamic command introspection.
func (c *Client) SetRootCommand(root *cobra.Command) {
	c.rootCmd = root
}

// GetSuggestion translates a natural-language request into a structured CommandSuggestion.
// Auto-detects Google Gemini vs Anthropic Claude based on API Key prefix.
func (c *Client) GetSuggestion(userInput, context string) (*CommandSuggestion, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("AI API key is not configured. Please add your Gemini or Anthropic API key in configuration")
	}

	userContent := strings.TrimSpace(userInput)
	if context != "" {
		userContent = fmt.Sprintf("Context: %s\n\nUser Request: %s", context, userContent)
	}

	systemPrompt := BuildSystemPrompt(c.rootCmd)

	// Auto-detect Provider: If starts with "sk-ant-", use Anthropic Claude; otherwise use Google Gemini
	if strings.HasPrefix(c.apiKey, "sk-ant-") {
		return c.callAnthropic(userContent, systemPrompt)
	}
	return c.callGemini(userContent, systemPrompt)
}

func (c *Client) callGemini(userContent, systemPrompt string) (*CommandSuggestion, error) {
	candidateModels := []string{"gemini-flash-latest", "gemini-flash-lite-latest", "gemini-2.5-flash-lite", "gemini-3.6-flash"}
	var lastErr error

	reqPayload := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: systemPrompt}},
		},
		Contents: []geminiContent{
			{
				Role:  "user",
				Parts: []geminiPart{{Text: userContent}},
			},
		},
		GenerationConfig: map[string]interface{}{
			"response_mime_type": "application/json",
		},
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode Gemini request: %w", err)
	}

	for _, model := range candidateModels {
		url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, c.apiKey)

		// Try up to 2 times for temporary 503 spikes
		for attempt := 0; attempt < 2; attempt++ {
			if attempt > 0 {
				time.Sleep(500 * time.Millisecond)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
			if err != nil {
				lastErr = err
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := c.httpClient.Do(req)
			if err != nil {
				lastErr = err
				continue
			}

			respBody, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				lastErr = err
				continue
			}

			var geminiResp geminiResponse
			if err := json.Unmarshal(respBody, &geminiResp); err != nil {
				lastErr = fmt.Errorf("failed to parse Gemini response: %w (raw: %s)", err, string(respBody))
				continue
			}

			if geminiResp.Error != nil {
				lastErr = fmt.Errorf("Gemini API Error (%d): %s", geminiResp.Error.Code, geminiResp.Error.Message)
				if geminiResp.Error.Code == 503 {
					continue // retry 503
				}
				break // try next model if 404/400
			}

			if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
				lastErr = fmt.Errorf("Gemini API returned empty response: %s", string(respBody))
				continue
			}

			rawText := CleanJSONOutput(geminiResp.Candidates[0].Content.Parts[0].Text)
			sug, err := parseSuggestionJSON(rawText)
			if err == nil {
				return sug, nil
			}
			lastErr = err
		}
	}

	return nil, lastErr
}

func (c *Client) callAnthropic(userContent, systemPrompt string) (*CommandSuggestion, error) {
	reqPayload := anthropicMessageRequest{
		Model:     c.model,
		MaxTokens: 1024,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userContent},
		},
	}

	payloadBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("POST", AnthropicAPIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", AnthropicVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr anthropicResponse
		if jsonErr := json.Unmarshal(respBody, &apiErr); jsonErr == nil && apiErr.Error != nil {
			return nil, fmt.Errorf("Anthropic API error (%d): %s", resp.StatusCode, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("Anthropic API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	var rawText string
	for _, block := range apiResp.Content {
		if block.Type == "text" || block.Text != "" {
			rawText += block.Text
		}
	}

	rawText = CleanJSONOutput(rawText)
	return parseSuggestionJSON(rawText)
}

func parseSuggestionJSON(rawText string) (*CommandSuggestion, error) {
	if strings.TrimSpace(rawText) == "" {
		return nil, fmt.Errorf("empty response received from AI model")
	}

	var suggestion CommandSuggestion
	if err := json.Unmarshal([]byte(rawText), &suggestion); err != nil {
		return nil, fmt.Errorf("failed to parse AI command suggestion JSON: %w (raw response: %s)", err, rawText)
	}

	if strings.TrimSpace(suggestion.Resource) == "" || strings.TrimSpace(suggestion.Subcommand) == "" {
		return nil, fmt.Errorf("invalid AI suggestion: missing resource or subcommand")
	}

	return &suggestion, nil
}

// CleanJSONOutput strips markdown code fences (```json ... ```) or whitespace from LLM output.
func CleanJSONOutput(raw string) string {
	cleaned := strings.TrimSpace(raw)

	// Strip ```json or ```
	if strings.HasPrefix(cleaned, "```") {
		idx := strings.Index(cleaned, "\n")
		if idx != -1 {
			cleaned = cleaned[idx+1:]
		} else {
			cleaned = strings.TrimPrefix(cleaned, "```json")
			cleaned = strings.TrimPrefix(cleaned, "```")
		}
	}

	if strings.HasSuffix(cleaned, "```") {
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
	}

	return strings.TrimSpace(cleaned)
}

// GetSuggestion is a package-level helper that uses the default configured client.
func GetSuggestion(userInput string, context string) (*CommandSuggestion, error) {
	client := NewClient()
	return client.GetSuggestion(userInput, context)
}
