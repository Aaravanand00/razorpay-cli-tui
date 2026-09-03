package ai

import (
	"fmt"
	"sort"
	"strings"
)

// CommandSuggestion represents an LLM's structured translation of a natural-language request into a CLI command.
type CommandSuggestion struct {
	Resource    string            `json:"resource"`
	Subcommand  string            `json:"subcommand"`
	Flags       map[string]string `json:"flags,omitempty"`
	Explanation string            `json:"explanation"`
	Confidence  string            `json:"confidence"` // "high" | "low"
}

// FullCommandString formats the suggestion into the full executable CLI command line.
func (c *CommandSuggestion) FullCommandString() string {
	parts := []string{"razorpay"}
	if c.Resource != "" {
		parts = append(parts, c.Resource)
	}
	if c.Subcommand != "" {
		parts = append(parts, c.Subcommand)
	}

	if len(c.Flags) > 0 {
		var flagKeys []string
		for k := range c.Flags {
			flagKeys = append(flagKeys, k)
		}
		sort.Strings(flagKeys)

		for _, k := range flagKeys {
			val := c.Flags[k]
			cleanKey := strings.TrimLeft(k, "-")
			if val == "" || val == "true" {
				parts = append(parts, fmt.Sprintf("--%s", cleanKey))
			} else {
				if strings.Contains(val, " ") {
					parts = append(parts, fmt.Sprintf("--%s \"%s\"", cleanKey, val))
				} else {
					parts = append(parts, fmt.Sprintf("--%s %s", cleanKey, val))
				}
			}
		}
	}

	return strings.Join(parts, " ")
}

// IsLowConfidence returns true if the suggestion is marked low confidence or requires extra review.
func (c *CommandSuggestion) IsLowConfidence() bool {
	return strings.ToLower(strings.TrimSpace(c.Confidence)) == "low"
}
