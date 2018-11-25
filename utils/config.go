package utils

import "github.com/spf13/viper"

// ReadConfig accepts a path, a filename and a map of default values.
// It returns an instance of Viper with loaded configurations:
// Notset < Default < Configfile < Environment < (flag)
func ReadConfig(path string, filename string, defaults map[string]interface{}) error {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigName(filename)
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
