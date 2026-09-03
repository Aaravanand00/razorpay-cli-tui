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
)

type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
	rootCmd    *cobra.Command
}

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

func NewClient() *Client {
	return NewClientWithKey(config.AIApiKey())
}

func NewClientWithKey(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		model:      DefaultClaudeModel,
		httpClient: &http.Client{Timeout: 20 * time.Second},
		rootCmd:    GetGlobalRootCommand(),
	}
}

// SetRootCommand sets the root command used for dynamic command introspection.
func (c *Client) SetRootCommand(root *cobra.Command) {
	c.rootCmd = root
}

// GetSuggestion translates a natural-language request into a structured CommandSuggestion using Anthropic API.
func (c *Client) GetSuggestion(userInput, context string) (*CommandSuggestion, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return nil, fmt.Errorf("AI API key is not configured. Please add your API key in configuration")
	}

	userContent := strings.TrimSpace(userInput)
	if context != "" {
		userContent = fmt.Sprintf("Context: %s\n\nUser Request: %s", context, userContent)
	}

	systemPrompt := BuildSystemPrompt(c.rootCmd)

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
		// Remove leading ```json or ```
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
