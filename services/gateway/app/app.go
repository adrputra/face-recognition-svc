package app

import (
	"face-recognition-svc/gateway/app/config"
	"face-recognition-svc/gateway/app/connection"
	"face-recognition-svc/gateway/app/model"
	"face-recognition-svc/gateway/app/router"
	"face-recognition-svc/gateway/app/utils"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	e := echo.New()

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: utils.LogError,
	}))

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	auth := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.Auth.AccessSecret),
	}

	public := e.Group("api")
	api := public.Group("/service")

	api.Use(echojwt.WithConfig(auth))
	api.Use(router.GetFactory().Middleware.Auth.IsAuthorized())

	e.Use(middleware.Logger())
	router.InitPublicRoute("", public)
	router.InitUserRoute("/user", api)
	router.InitDatasetRoute("/dataset", api)
	router.InitRoleRoute("/role", api)
	router.InitPermissionRoute("/permission", api)
	router.InitFeatureRoute("/feature", api)
	router.InitParamRoute("/param", api)
	router.InitInstitutionRoute("/institution", api)

	e.Logger.Fatal(e.Start(host + ":" + strconv.Itoa(port)))
}
