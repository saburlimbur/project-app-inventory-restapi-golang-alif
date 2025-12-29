package utils

import (
	"errors"

	"github.com/spf13/viper"
)

type Configuration struct {
	AppName  string
	Port     string
	PathLogg string
}

type DatabaseCofig struct {
	Name     string
	Username string
	Password string
	Host     string
	Port     string
	MaxConn  int32
}

func ReadConfigurationEnv() (Configuration, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, errors.New("err loaded env files")
	}

	viper.AutomaticEnv()

	return Configuration{
		AppName:  viper.GetString("APP_NAME"),
		Port:     viper.GetString("PORT"),
		PathLogg: viper.GetString("PATH_LOGGING"),
	}, nil
}
