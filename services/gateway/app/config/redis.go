package config

type Redis struct {
	Host     string `yaml:"host" default:"127.0.0.1" desc:"config:sql:host"`
	Port     int    `yaml:"port" default:"3306" desc:"config:sql:port"`
	Username string `yaml:"username" default:"root"  desc:"config:sql:username"`
	Password string `yaml:"password" default:"root" desc:"config:sql:password"`
}
