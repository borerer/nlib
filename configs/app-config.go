package configs

import (
	"strings"

	"github.com/borerer/nlib/utils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	LogLevel string     `json:"log_level" mapstructure:"log_level"`
	Addr     string     `json:"addr" mapstructure:"addr"`
	Port     string     `json:"port" mapstructure:"port"`
	Mongo    string     `json:"mongo" mapstructure:"mongo"`
	File     FileConfig `json:"file" mapstructure:"file"`
}

type FileConfig struct {
	Engine string      `json:"engine" mapstructure:"engine"`
	Minio  MinioConfig `json:"minio" mapstructure:"minio"`
	FS     FSConfig    `json:"fs" mapstructure:"fs"`
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint"`
	AccessKey string `json:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key"`
	UseSSL    bool   `json:"use_ssl" mapstructure:"use_ssl"`
	Bucket    string `json:"bucket" mapstructure:"bucket"`
}

type FSConfig struct {
	Dir string `json:"dir" mapstructure:"dir"`
}

func GetAppConfig() *AppConfig {
	viper.AddConfigPath("data")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetEnvPrefix("nlib")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "9502")
	utils.Must(viper.ReadInConfig())
	var appConfig AppConfig
	utils.Must(viper.Unmarshal(&appConfig))
	return &appConfig
}
