package util

import (
	"github.com/spf13/viper"
	"time"
)

// contain all of the variable for the environemnt
type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (conf Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") //json //xml
	viper.AutomaticEnv()       //json //xml

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&conf)

	return

}
