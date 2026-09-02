package components

import (
	"strings"

	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func RenderToast(s *state.SessionState) string {
	if s.Toast == nil || s.Toast.Message == "" {
		return ""
	}

	var content string
	if s.Toast.IsError {
		content = styles.ToastErrorStyle.Render("✖  " + s.Toast.Message)
	} else if strings.Contains(s.Toast.Message, "▲") || strings.Contains(s.Toast.Message, "Test") {
		content = styles.ToastWarningStyle.Render("▲  " + s.Toast.Message)
	} else {
		content = styles.ToastSuccessStyle.Render("✔  " + s.Toast.Message)
	}

	// Clean vertical breathing margin around toast for maximum clarity
	return "\n" + content + "\n"
}
