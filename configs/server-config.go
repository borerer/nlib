package configs

import (
	"strings"

	"github.com/borerer/nlib/utils"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	LogLevel string    `yaml:"log-level" mapstructure:"log-level"`
	API      APIConfig `yaml:"api" mapstructure:"api"`
}

type APIConfig struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port string `yaml:"port" mapstructure:"port"`
}

type BuiltinConfig struct {
	KV    KVConfig `yaml:"kv" mapstructure:"kv"`
	Logs  KVConfig `yaml:"logs" mapstructure:"logs"`
	Files KVConfig `yaml:"files" mapstructure:"files"`
}

type KVConfig struct {
	Enabled bool        `yaml:"enabled" mapstructure:"enabled"`
	Mongo   MongoConfig `yaml:"mongo" mapstructure:"mongo"`
}

type LogsConfig struct {
	Enabled bool        `yaml:"enabled" mapstructure:"enabled"`
	Mongo   MongoConfig `yaml:"mongo" mapstructure:"mongo"`
}

type FilesConfig struct {
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
}

type MongoConfig struct {
	URI        string `yaml:"uri" mapstructure:"uri"`
	Database   string `yaml:"database" mapstructure:"database"`
	Collection string `yaml:"collection" mapstructure:"collection"`
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
