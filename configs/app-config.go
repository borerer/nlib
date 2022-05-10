package configs

import (
	"gitea.home.iloahz.com/iloahz/nlib/utils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	LogLevel string `json:"log_level" mapstructure:"log_level"`
}

var (
	appConfig = &AppConfig{}
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	utils.Must(viper.ReadInConfig())
	utils.Must(viper.Unmarshal(appConfig))
}

func GetAppConfig() *AppConfig {
	return appConfig
}
