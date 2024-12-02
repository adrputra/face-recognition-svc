package config

type MinioS3 struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	SecretKey string `yaml:"secretKey"`
	Tls       bool   `yaml:"tls"`
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
}
