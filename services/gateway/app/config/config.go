package config

type Config struct {
	Listener struct {
		Host string
		Port int
	}
	DatabaseProfile struct {
		Database Database `yaml:"database"`
		Timeout  int      `yaml:"timeout" default:"30000"`
	} `yaml:"databaseProfile"`
	Auth         Auth        `yaml:"auth"`
	Redis        Redis       `yaml:"redis"`
	Jaeger       Jaeger      `yaml:"jaeger"`
	MinioProfile MinioS3     `yaml:"minioProfile"`
	API          APIEndpoint `yaml:"api"`
	RabbitMQ     RabbitMQ    `yaml:"rabbitmq"`
}

var config *Config

func GetConfig() *Config {
	return config
}
