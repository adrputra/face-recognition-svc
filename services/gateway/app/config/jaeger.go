package config

type Jaeger struct {
	Host           string `yaml:"host" default:"127.0.0.1" desc:"config:jaeger:host"`
	Port           string `yaml:"port" default:"3306" desc:"config:jaeger:port"`
	TracePerSecond int    `yaml:"tracePerSecond" default:"100" desc:"config:jaeger:tracePerSecond"`
	ServiceName    string `yaml:"serviceName" default:"mrscustom" desc:"config:jaeger:serviceName"`
}
