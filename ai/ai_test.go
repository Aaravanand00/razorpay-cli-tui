package ai

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCommandSuggestionFullCommandString(t *testing.T) {
	s := CommandSuggestion{
		Resource:   "payments",
		Subcommand: "list",
		Flags: map[string]string{
			"status": "failed",
			"count":  "10",
		},
		Explanation: "List 10 failed payments",
		Confidence:  "high",
	}

	cmdStr := s.FullCommandString()
	if !strings.Contains(cmdStr, "razorpay payments list") {
		t.Fatalf("expected command string to contain 'razorpay payments list', got %s", cmdStr)
	}
	if !strings.Contains(cmdStr, "--status failed") || !strings.Contains(cmdStr, "--count 10") {
		t.Fatalf("expected flags in command string, got %s", cmdStr)
	}
	if s.IsLowConfidence() {
		t.Fatal("expected high confidence")
	}
}

func TestIntrospectCommandsAndPrompt(t *testing.T) {
	root := &cobra.Command{Use: "razorpay"}
	ordersCmd := &cobra.Command{Use: "orders", Short: "Manage orders"}
	createCmd := &cobra.Command{Use: "create", Short: "Create a new order"}
	createCmd.Flags().Int("amount", 0, "Amount in paise")
	ordersCmd.AddCommand(createCmd)
	root.AddCommand(ordersCmd)

	catalog := IntrospectCommands(root)
	if !strings.Contains(catalog, "razorpay orders create: Create a new order") {
		t.Fatalf("expected catalog to contain 'razorpay orders create', got:\n%s", catalog)
	}
	if !strings.Contains(catalog, "Flags: --amount") {
		t.Fatalf("expected catalog to contain flag description, got:\n%s", catalog)
	}

	prompt := BuildSystemPrompt(root)
	if !strings.Contains(prompt, "AVAILABLE RAZORPAY CLI COMMANDS:") {
		t.Fatal("expected system prompt header")
	}
	if !strings.Contains(prompt, "RESPONSE SCHEMA:") {
		t.Fatal("expected response schema in system prompt")
	}
}

func TestCleanJSONOutput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "```json\n{\"resource\":\"payments\",\"subcommand\":\"list\"}\n```",
			expected: "{\"resource\":\"payments\",\"subcommand\":\"list\"}",
		},
		{
			input:    "```\n{\"resource\":\"orders\",\"subcommand\":\"create\"}\n```",
			expected: "{\"resource\":\"orders\",\"subcommand\":\"create\"}",
		},
		{
			input:    "  {\"resource\":\"refunds\",\"subcommand\":\"create\"}  ",
			expected: "{\"resource\":\"refunds\",\"subcommand\":\"create\"}",
		},
	}

	for _, c := range cases {
		cleaned := CleanJSONOutput(c.input)
		if cleaned != c.expected {
			t.Fatalf("expected %s, got %s", c.expected, cleaned)
		}
		var target CommandSuggestion
		if err := json.Unmarshal([]byte(cleaned), &target); err != nil {
			t.Fatalf("failed to unmarshal cleaned JSON: %v", err)
		}
	}
}
