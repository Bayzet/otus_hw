package main

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
	HTTP    HTTPConfig  `yaml:"http"`
	GRPC    GRPCConfig  `yaml:"grpc"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type StorageConf struct {
	Type   string `yaml:"type"`
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
