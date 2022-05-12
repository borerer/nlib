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

func GetAppConfig() *AppConfig {
	viper.AddConfigPath("data")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.SetEnvPrefix("nlib")
	viper.BindEnv("addr") // will bind NLIB_ADDR
	viper.BindEnv("port")
	viper.SetDefault("addr", "0.0.0.0")
	viper.SetDefault("port", "8080")
	utils.Must(viper.ReadInConfig())
	var appConfig AppConfig
	utils.Must(viper.Unmarshal(&appConfig))
	return &appConfig
}
