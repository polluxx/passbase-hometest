package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const configName = "config"

type Config struct {
	Database
	Rates
}

type Database struct {
	Type     string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Rates struct {
	Fixer
}

type Fixer struct {
	APIKey   string
	BaseHost string
}

func LoadFromPath(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}

		if path != "" {
			return nil, fmt.Errorf("cannot find config path: %w", err)
		}
	}

	var conf Config
	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
