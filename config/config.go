package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razorpay/razorpay-cli/output"
	"github.com/spf13/viper"
)

const (
	configDir  = ".razorpay"
	configFile = "config"
	configType = "yaml"
)

func ConfigFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, configDir, configFile+"."+configType)
}

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error finding home directory: %v\n", err)
		os.Exit(1)
	}

	viper.SetConfigName(configFile)
	viper.SetConfigType(configType)
	viper.AddConfigPath(filepath.Join(home, configDir))
	viper.AutomaticEnv()

	// env overrides
	viper.SetEnvPrefix("RAZORPAY")
	_ = viper.BindEnv("key_id", "RAZORPAY_KEY_ID")
	_ = viper.BindEnv("key_secret", "RAZORPAY_KEY_SECRET")
	_ = viper.BindEnv("output_format", "RAZORPAY_OUTPUT_FORMAT")

	viper.SetDefault("output_format", output.DefaultFormat)
	viper.SetDefault("active_mode", "test")

	_ = viper.ReadInConfig()
}

func KeyID() string {
	mode := ActiveMode()
	if mode == "live" && LiveKeyID() != "" {
		return LiveKeyID()
	}
	if mode == "test" && TestKeyID() != "" {
		return TestKeyID()
	}
	k := viper.GetString("key_id")
	if strings.HasPrefix(k, "rzp_test_") || strings.HasPrefix(k, "rzp_live_") {
		return k
	}
	return ""
}

func KeySecret() string {
	mode := ActiveMode()
	if mode == "live" && LiveKeySecret() != "" {
		return LiveKeySecret()
	}
	if mode == "test" && TestKeySecret() != "" {
		return TestKeySecret()
	}
	return viper.GetString("key_secret")
}

func ActiveMode() string {
	m := strings.ToLower(strings.TrimSpace(viper.GetString("active_mode")))
	if m == "" {
		if strings.HasPrefix(viper.GetString("key_id"), "rzp_live_") {
			return "live"
		}
		return "test"
	}
	return m
}

func SetActiveMode(mode string) {
	viper.Set("active_mode", strings.ToLower(strings.TrimSpace(mode)))
}

func TestKeyID() string {
	k := viper.GetString("test.key_id")
	if strings.HasPrefix(k, "rzp_test_") {
		return k
	}
	return ""
}

func TestKeySecret() string {
	if TestKeyID() == "" {
		return ""
	}
	return viper.GetString("test.key_secret")
}

func LiveKeyID() string {
	k := viper.GetString("live.key_id")
	if strings.HasPrefix(k, "rzp_live_") {
		return k
	}
	return ""
}

func LiveKeySecret() string {
	if LiveKeyID() == "" {
		return ""
	}
	return viper.GetString("live.key_secret")
}

// OutputFormat returns the configured presentation format (json, yaml, …).
func OutputFormat() string {
	return strings.ToLower(strings.TrimSpace(viper.GetString("output_format")))
}

// SetOutputFormat updates the in-memory output_format value.
func SetOutputFormat(format string) {
	viper.Set("output_format", strings.ToLower(strings.TrimSpace(format)))
}

func Save(keyID, keySecret string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	viper.Set("key_id", keyID)
	viper.Set("key_secret", keySecret)

	if strings.HasPrefix(keyID, "rzp_live_") {
		viper.Set("active_mode", "live")
		viper.Set("live.key_id", keyID)
		viper.Set("live.key_secret", keySecret)
	} else if strings.HasPrefix(keyID, "rzp_test_") {
		viper.Set("active_mode", "test")
		viper.Set("test.key_id", keyID)
		viper.Set("test.key_secret", keySecret)
	}

	return viper.WriteConfigAs(filepath.Join(dir, configFile+"."+configType))
}

func SaveDualConfig(activeMode, testKeyID, testKeySecret, liveKeyID, liveKeySecret string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, configDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	viper.Set("active_mode", activeMode)
	viper.Set("test.key_id", testKeyID)
	viper.Set("test.key_secret", testKeySecret)
	viper.Set("live.key_id", liveKeyID)
	viper.Set("live.key_secret", liveKeySecret)

	// Set active primary keys for backward compatibility
	if activeMode == "live" && strings.HasPrefix(liveKeyID, "rzp_live_") {
		viper.Set("key_id", liveKeyID)
		viper.Set("key_secret", liveKeySecret)
	} else if strings.HasPrefix(testKeyID, "rzp_test_") {
		viper.Set("key_id", testKeyID)
		viper.Set("key_secret", testKeySecret)
	}

	return viper.WriteConfigAs(filepath.Join(dir, configFile+"."+configType))
}
