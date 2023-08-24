package config

import (
	"strconv"
	"github.com/spf13/viper"
	"github.com/hope-ag/go-dynamo/utils/env"
)

type Config struct {
	Port, Timeout int
	Dialect, DatabaseURI string
}

func GetConfig() Config {
	return Config{
		Port: viper.GetInt("PORT"),
		Timeout: viper.GetInt("TIMEOUT"),
		Dialect: viper.GetString("DB_DIALECT"),
		DatabaseURI: viper.GetString("DATABASE_URI"),
	}
}