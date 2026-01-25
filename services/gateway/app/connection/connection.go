package connection

import (
	"context"
	"crypto/tls"
	"face-recognition-svc/gateway/app/config"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db      *gorm.DB
	Storage *s3.S3
	Redis   *redis.Client
	Mq      *amqp.Channel
)

func InitConnection(c config.Config) {
	Db = NewDatabaseConnection(&c.DatabaseProfile.Database)
	Storage = NewStorageConnection(&c.MinioProfile)
	Redis = NewRedisConnection(&c.Redis, context.Background())
	Mq = NewRabbitMQConnection(&c.RabbitMQ)
}

func NewDatabaseConnection(c *config.Database) *gorm.DB {
	// Use UTC in DSN to avoid DB time zone lookup errors.
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		c.Host,
		c.Username,
		c.Password,
		c.Database,
		c.Port)

	log.Info().Msg(dataSourceName)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		log.Panic().Str("database", c.Database).Err(err).Msg("Cannot Connect To Database")
	}

	log.Info().Str("database", c.Database).Msg("Connected To Database")

	return db
}

func NewStorageConnection(cfg *config.MinioS3) *s3.S3 {
	awsAccessKey := cfg.Username
	awsSecretKey := cfg.SecretKey

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Bypass certificate verification
		},
	}

	// Create a custom HTTP client with the custom transport
	httpClient := &http.Client{
		Transport: transport,
	}

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			awsAccessKey,
			awsSecretKey,
			"",
		),
		Endpoint:         aws.String(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)),
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(cfg.Region),
		HTTPClient:       httpClient,
	})

	if err != nil {
		log.Panic().Err(err).Msg("Cannot Connect To Minio")
	}

	log.Info().Str("username", cfg.Username).Str("host", cfg.Host).Str("port", cfg.Port).Msg("Minio credentials")
	log.Info().Str("host", cfg.Host).Str("port", cfg.Port).Msg("Connected To Minio")

	return s3.New(sess)
}

func NewRedisConnection(c *config.Redis, ctx context.Context) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port), // Replace with your Redis server address
		Password: c.Password,                           // No password set (use if your Redis requires auth)
		DB:       0,                                    // Default DB
	})

	// Test the connection
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect to Redis")
	}
	log.Info().Str("pong", pong).Msg("Connected to Redis")

	return rdb

}

func NewRabbitMQConnection(c *config.RabbitMQ) *amqp.Channel {
	log.Info().Str("url", fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port)).Msg("RabbitMQ connection string")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port))
	if err != nil {
		log.Panic().Err(err).Msg("Cannot Connect To RabbitMQ")
	}
	log.Info().Str("host", c.Host).Str("port", c.Port).Msg("Connected To RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open a channel")
	}

	return ch
}
