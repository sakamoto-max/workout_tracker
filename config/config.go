package config

import (
	"log"
	"github.com/spf13/viper"
)

//
// DB_NAME=WORKOUT
// DB_PASSWORD=root
// DB_OWNER_NAME=postgres

// WEB_PORT=5000

var Config configStruct

type configStruct struct {
	DbPort string `mapstructure:"DB_PORT"`
	Db string `mapstructure:"DB"`
	DbName string `mapstructure:"DB_NAME"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbOwnerName string `mapstructure:"DB_OWNER_NAME"`
	WebPort string `mapstructure:"WEB_PORT"`
	DbHost string `mapstructure:"DB_HOST"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func InitializeConfig() {
	data := LoadValuesFromEnv()

	Config = data
}

func LoadValuesFromEnv() (configStruct) {

	var config configStruct

	viper.SetConfigName("app")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error while loading the .env file : %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error while unmarshalling the env variables : %v", err)
	}

	return config
}