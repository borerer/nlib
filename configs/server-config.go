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
	Mongo  string       `yaml:"mongo" mapstructure:"mongo"`
	Webdav WebdavConfig `yaml:"webdav" mapstructure:"webdav"`
	Samba  SambaConfig  `yaml:"samba" mapstructure:"samba"`
	Minio  MinioConfig  `yaml:"minio" mapstructure:"minio"`
}

type WebdavConfig struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
}

type SambaConfig struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	Share    string `yaml:"share" mapstructure:"share"`
	Path     string `yaml:"path" mapstructure:"path"`
}

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKey string `yaml:"access-key" mapstructure:"access-key"`
	SecretKey string `yaml:"secret-key" mapstructure:"secret-key"`
	UseSSL    bool   `yaml:"use-ssl" mapstructure:"use-ssl"`
	Bucket    string `yaml:"bucket" mapstructure:"bucket"`
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
