package config

type Auth struct {
	AccessSecret  string `yaml:"accessSecret"`
	RefreshSecret string `yaml:"refreshSecret"`
	AccessExpiry  string `yaml:"accessExpiry"`
}
