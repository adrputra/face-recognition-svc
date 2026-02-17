package config

type APIEndpoint struct {
	ProcessingSVC struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Endpoint string `yaml:"endpoint"`
	} `yaml:"processingsvc"`
	PresenceSVC struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"presencesvc"`
}
