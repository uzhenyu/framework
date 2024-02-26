package config

import "github.com/spf13/viper"

func ReadConfig(filePath string) {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		return
	}
}
