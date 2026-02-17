package app

import (
	"os"
	"strconv"

	"github.com/adrputra/face-recognition-svc/gateway-gin/app/config"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/connection"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/router"
	"github.com/adrputra/face-recognition-svc/gateway-gin/app/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Start() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize console writer, disable on production
	if os.Getenv("ENV") == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	config.InitConfig()
	cfg := config.GetConfig()

	utils.InitTimeLocation()

	tracer, closer, err := utils.InitJaeger(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Jaeger tracer")
	}
	defer closer.Close()

	// Set global tracer
	opentracing.SetGlobalTracer(tracer)

	connection.InitConnection(*cfg)
	connection.MigrateDatabase(&cfg.DatabaseProfile.Database)
	router.InitFactory(cfg, connection.Db, connection.Storage, connection.Redis, connection.Mq)

	host := cfg.Listener.Host
	port := cfg.Listener.Port

	e := gin.New()
	e.Use(gin.Recovery())
	e.Use(gin.Logger())

	e.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	public := e.Group("/api")
	api := public.Group("/service")

	api.Use(router.JWTMiddleware(cfg.Auth.AccessSecret))
	router.InitPublicRoute("", public)
	router.InitUserRoute("/user", api)
	router.InitDatasetRoute("/dataset", api)
	router.InitRoleRoute("/role", api)
	router.InitPermissionRoute("/permission", api)
	router.InitFeatureRoute("/feature", api)
	router.InitParamRoute("/param", api)
	router.InitInstitutionRoute("/institution", api)

	if err := e.Run(host + ":" + strconv.Itoa(port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to start gateway-gin service")
	}
}
