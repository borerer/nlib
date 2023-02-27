package configs

import (
	"strings"

	"github.com/borerer/nlib/utils"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	LogLevel string        `yaml:"log-level" mapstructure:"log-level"`
	API      APIConfig     `yaml:"api" mapstructure:"api"`
	Builtin  BuiltinConfig `yaml:"builtin" mapstructure:"builtin"`
}

type APIConfig struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port string `yaml:"port" mapstructure:"port"`
}

type BuiltinConfig struct {
	Echo  EchoConfig  `yaml:"echo" mapstructure:"echo"`
	KV    KVConfig    `yaml:"kv" mapstructure:"kv"`
	Logs  LogsConfig  `yaml:"logs" mapstructure:"logs"`
	Files FilesConfig `yaml:"files" mapstructure:"files"`
}

type EchoConfig struct {
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`
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
	Enabled bool         `yaml:"enabled" mapstructure:"enabled"`
	Backend string       `yaml:"backend" mapstructure:"backend"`
	Minio   MinioConfig  `yaml:"minio" mapstructure:"minio"`
	Webdav  WebdavConfig `yaml:"webdav" mapstructure:"webdav"`
}

type MongoConfig struct {
	URI        string `yaml:"uri" mapstructure:"uri"`
	Database   string `yaml:"database" mapstructure:"database"`
	Collection string `yaml:"collection" mapstructure:"collection"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKey string `yaml:"access-key" mapstructure:"access-key"`
	SecretKey string `yaml:"secret-key" mapstructure:"secret-key"`
	UseSSL    bool   `yaml:"use-ssl" mapstructure:"use-ssl"`
	Bucket    string `yaml:"bucket" mapstructure:"bucket"`
}

type WebdavConfig struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

func GetServerConfig() *ServerConfig {
	// config.yaml
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// env vars
	viper.SetEnvPrefix("nlib")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// read config
	utils.Must(viper.ReadInConfig())
	var serverConfig ServerConfig
	utils.Must(viper.Unmarshal(&serverConfig))
	return &serverConfig
}
