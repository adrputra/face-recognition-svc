package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InitConfig() {
	basedir := filepath.Join(".")
	viper.AddConfigPath(basedir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config.yaml")

	if err := viper.MergeInConfig(); err != nil {
		log.Panic().Err(err).Msg("Failed to load config")
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)

		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			expr := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")

			switch {
			case strings.HasPrefix(expr, "file:"):
				secretPath := strings.TrimPrefix(expr, "file:")
				secret, err := os.ReadFile(secretPath)
				if err != nil {
					log.Panic().
						Err(err).
						Str("path", secretPath).
						Msg("Failed to read secret file")
				}
				viper.Set(k, strings.TrimSpace(string(secret)))

			default:
				viper.Set(k, getEnvOrPanic(expr))
			}
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Panic().Err(err).Msg("Failed to unmarshal config")
	}
}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(env) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
