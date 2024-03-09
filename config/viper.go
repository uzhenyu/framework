package config

import "github.com/spf13/viper"

func ReadConfig(filePath string) error {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
