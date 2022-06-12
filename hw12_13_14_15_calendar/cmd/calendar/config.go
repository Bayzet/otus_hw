package main

type Config struct {
	Logger LoggerConf `yaml:"logger"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func NewConfig() Config {
	return Config{}
}
