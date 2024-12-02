package connection

import (
	"context"
	"crypto/tls"
	"face-recognition-svc/app/config"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
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
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FJakarta",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logrus.Panicf("Cannot Connect To Database %s: %v", c.Database, err)
	}

	logrus.Printf("Connected To Database %s", c.Database)

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
		logrus.Panicf("Cannot Connect To Minio: %v", err)
	}

	logrus.Println(cfg.Username, cfg.SecretKey)
	logrus.Printf("Connected To Minio at %s:%s", cfg.Host, cfg.Port)

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
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

	return rdb

}

func NewRabbitMQConnection(c *config.RabbitMQ) *amqp.Channel {
	logrus.Printf(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port))
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", c.Username, c.Password, c.Host, c.Port))
	if err != nil {
		logrus.Panicf("Cannot Connect To RabbitMQ: %v", err)
	}
	logrus.Printf("Connected To RabbitMQ at %s:%s", c.Host, c.Port)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	return ch
}
