package components

import (
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func RenderToast(s *state.SessionState) string {
	if s.Toast == nil || s.Toast.Message == "" {
		return ""
	}

	if s.Toast.IsError {
		return styles.ToastErrorStyle.Render("✖  " + s.Toast.Message)
	}
	return styles.ToastSuccessStyle.Render("✔  " + s.Toast.Message)
}
