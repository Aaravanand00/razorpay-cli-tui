package ai

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var globalRootCmd *cobra.Command

// SetGlobalRootCommand sets the global root cobra command for dynamic introspection.
func SetGlobalRootCommand(root *cobra.Command) {
	globalRootCmd = root
}

// GetGlobalRootCommand returns the registered global root cobra command.
func GetGlobalRootCommand() *cobra.Command {
	return globalRootCmd
}

// IntrospectCommands dynamically builds a formatted catalog of all registered Cobra resources, subcommands, and flags.
func IntrospectCommands(root *cobra.Command) string {
	if root == nil {
		root = globalRootCmd
	}
	if root == nil {
		return ""
	}

	var b strings.Builder
	for _, cmd := range root.Commands() {
		if cmd.Hidden || cmd.Name() == "help" || cmd.Name() == "completion" {
			continue
		}
		formatCommandTree(&b, cmd, "")
	}
	return b.String()
}

func formatCommandTree(b *strings.Builder, cmd *cobra.Command, parentPath string) {
	currentPath := cmd.Name()
	if parentPath != "" {
		currentPath = parentPath + " " + cmd.Name()
	}

	subCommands := cmd.Commands()
	var visibleSubs []*cobra.Command
	for _, sc := range subCommands {
		if !sc.Hidden && sc.Name() != "help" {
			visibleSubs = append(visibleSubs, sc)
		}
	}

	if len(visibleSubs) > 0 {
		for _, sc := range visibleSubs {
			formatCommandTree(b, sc, currentPath)
		}
		return
	}

	// Leaf command
	short := strings.TrimSpace(cmd.Short)
	if short == "" {
		short = strings.TrimSpace(cmd.Use)
	}

	b.WriteString(fmt.Sprintf("- razorpay %s: %s\n", currentPath, short))

	// Collect flags
	var flags []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name != "help" && f.Name != "version" {
			flags = append(flags, fmt.Sprintf("--%s (%s)", f.Name, f.Usage))
		}
	})

	if len(flags) > 0 {
		b.WriteString(fmt.Sprintf("    Flags: %s\n", strings.Join(flags, ", ")))
	}
}

// BuildSystemPrompt constructs the strict JSON system prompt with dynamic command catalog and safety rules.
func BuildSystemPrompt(root *cobra.Command) string {
	commandCatalog := IntrospectCommands(root)

	return fmt.Sprintf(`You are the Razorpay CLI AI Assistant.
Your job is to translate a user's natural language request into an exact, valid Razorpay CLI command and flags.

### AVAILABLE RAZORPAY CLI COMMANDS:
%s

### RESPONSE SCHEMA:
You must respond with ONLY a single valid JSON object matching this schema (NO markdown formatting, NO code blocks, NO backticks, NO extra prose):
{
  "resource": "<resource_name>",
  "subcommand": "<subcommand_name>",
  "flags": {
    "<flag_name>": "<flag_value>"
  },
  "explanation": "<short 1-sentence explanation of what this command will do>",
  "confidence": "<high|low>"
}

### STRICT RULES:
1. NEVER invent any resource or subcommand that is not in the list above.
2. If the user refers to an amount in Rupees (₹), convert it to paise (multiply by 100) if the flag expects paise (e.g. ₹500 -> 50000).
3. If the user refers to relative dates like "yesterday", "last week", calculate or set the appropriate epoch timestamps or date strings supported by the flag.
4. Mark "confidence" as "low" if:
   - The user request is ambiguous, vague, or missing required identifiers.
   - The action is destructive or modifies data (e.g., refund, update, capture, cancel, delete, contest).
5. Mark "confidence" as "high" ONLY for unambiguous read-only operations (e.g., list, fetch, card details).
6. Output ONLY raw JSON. Do not wrap in `+"`"+`json ... `+"`"+` code fences.`, commandCatalog)
}
