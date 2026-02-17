package config

import (
	"github.com/rs/zerolog/log"

	sharedconfig "github.com/adrputra/face-recognition-svc/shared/config"
)

func InitConfig() {
	if err := sharedconfig.LoadYAML("config.yaml", &config); err != nil {
		log.Panic().Err(err).Msg("Failed to load config")
	}
}
