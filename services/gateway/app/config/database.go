package config

type Database struct {
	Host     string `yaml:"host" default:"127.0.0.1" desc:"config:sql:host"`
	Port     int    `yaml:"port" default:"3306" desc:"config:sql:port"`
	Username string `yaml:"username" default:"root"  desc:"config:sql:username"`
	Password string `yaml:"password" default:"root" desc:"config:sql:password"`
	Database string `yaml:"database" default:"mydb" desc:"config:sql:database"`
	MaxConn  int    `yaml:"maxConn" default:"5"  desc:"config:sql:maxConn"`
	MaxIdle  int    `yaml:"maxIdle" default:"5"  desc:"config:sql:maxIdle"`
	LifeTime int    `yaml:"lifeTime" default:"5"  desc:"config:sql:lifeTime"`
}