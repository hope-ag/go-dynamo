package config

import "github.com/spf13/viper"

type Config struct {
	Port, Timeout        int
	Dialect, DatabaseURI string
}

func GetConfig() Config {
	return Config{
		Port:        viper.GetInt("PORT"),
		Timeout:     viper.GetInt("TIMEOUT"),
		Dialect:     viper.GetString("DB_DIALECT"),
		DatabaseURI: viper.GetString("DATABASE_URI"),
	}
}
