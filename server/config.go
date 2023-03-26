package server

import "github.com/spf13/viper"

type Config struct {
	v *viper.Viper
}

func NewConfig(v *viper.Viper) *Config {
	return &Config{v: v}
}
