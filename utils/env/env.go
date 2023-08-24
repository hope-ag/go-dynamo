package env

import (
	"os"

	"github.com/spf13/viper"
)

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func InitEnvironment() error {
	isProd := IsProduction()

	if !isProd {
		viper.SetConfigFile(".env")
		// return
	} else {
		viper.SetConfigFile(".env.production")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}