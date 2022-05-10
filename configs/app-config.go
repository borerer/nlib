package configs

import (
	"gitea.home.iloahz.com/iloahz/nlib/utils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	LogLevel string `json:"log_level" mapstructure:"log_level"`
	Addr     string `json:"addr" mapstructure:"addr"`
	Port     string `json:"port" mapstructure:"port"`
}

var (
	appConfig = &AppConfig{}
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "8080")
	utils.Must(viper.ReadInConfig())
	utils.Must(viper.Unmarshal(appConfig))
}

func GetAppConfig() *AppConfig {
	return appConfig
}
