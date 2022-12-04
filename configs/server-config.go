package configs

import (
	"strings"

	"github.com/borerer/nlib/utils"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	LogLevel string      `yaml:"log-level" mapstructure:"log-level"`
	API      APIConfig   `yaml:"api" mapstructure:"api"`
	Minio    MinioConfig `yaml:"minio" mapstructure:"minio"`
	Mongo    MongoConfig `yaml:"mongo" mapstructure:"mongo"`
}

type APIConfig struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port string `yaml:"port" mapstructure:"port"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKey string `yaml:"access-key" mapstructure:"access-key"`
	SecretKey string `yaml:"secret-key" mapstructure:"secret-key"`
	UseSSL    bool   `yaml:"use-ssl" mapstructure:"use-ssl"`
	Bucket    string `yaml:"bucket" mapstructure:"bucket"`
}

type MongoConfig struct {
	URI      string `yaml:"uri" mapstructure:"uri"`
	Database string `yaml:"database" mapstructure:"database"`
}

func GetServerConfig() *ServerConfig {
	// config.yaml
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// env vars
	// export NLIB_LOG_LEVEL=debug
	// export NLIB_API_PORT=9502
	// export NLIB_API_ADDR=0.0.0.0
	// export NLIB_MINIO_ENDPOINT=
	// export NLIB_MINIO_ACCESS_KEY=
	// export NLIB_MINIO_SECRET_KEY=
	// export NLIB_MINIO_USE_SSL=
	// export NLIB_MINIO_BUCKET=
	// export NLIB_MONGO_URI=
	// export NLIB_MONGO_DATABASE=
	viper.SetEnvPrefix("nlib")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// read config
	utils.Must(viper.ReadInConfig())
	var serverConfig ServerConfig
	utils.Must(viper.Unmarshal(&serverConfig))
	return &serverConfig
}
