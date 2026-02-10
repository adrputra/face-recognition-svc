package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadYAML loads a YAML file and expands values in the form:
// - ${ENV_VAR}
// - ${file:/path/to/secret}
func LoadYAML(path string, target interface{}) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("config path is required")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	for _, key := range v.AllKeys() {
		value := v.GetString(key)
		if isTemplateValue(value) {
			resolved, err := resolveTemplateValue(value)
			if err != nil {
				return err
			}
			v.Set(key, resolved)
		}
	}

	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}

func isTemplateValue(value string) bool {
	return strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}")
}

func resolveTemplateValue(value string) (string, error) {
	expr := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
	switch {
	case strings.HasPrefix(expr, "file:"):
		secretPath := strings.TrimPrefix(expr, "file:")
		secret, err := os.ReadFile(secretPath)
		if err != nil {
			return "", fmt.Errorf("read secret file %s: %w", secretPath, err)
		}
		return strings.TrimSpace(string(secret)), nil
	default:
		envValue := os.Getenv(expr)
		if envValue == "" {
			return "", fmt.Errorf("mandatory env variable not found: %s", expr)
		}
		return envValue, nil
	}
}
