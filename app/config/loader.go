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
	err := viper.MergeInConfig()

	if err != nil {
		log.Panic().Err(err).Msg("Failed to load config")
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			viper.Set(k, getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")))
		}
	}

	viper.Unmarshal(&config)

}

func getEnvOrPanic(env string) string {
	res := os.Getenv(env)
	if len(env) == 0 {
		panic("Mandatory env variable not found:" + env)
	}
	return res
}
