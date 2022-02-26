package config

import (
	"github.com/spf13/viper"
	"os"
	"time"
)

type (
	Config struct {
		HTTP  HTTPConfig
		PGSQL PGSQLConfig
	}

	HTTPConfig struct {
		Host               string        `mapstructure:"host"`
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	PGSQLConfig struct {
		Host         string
		User         string
		Password     string
		DatabaseName string `mapstructure:"dbname"`
		SSLMode      string `mapstructure:"sslmode"`
		Port         string `mapstructure:"port"`
	}
)

func Init(configPath string) (*Config, error) {
	if err := parseConfigFile(configPath); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setEnv(&cfg)

	return &cfg, nil
}

func setEnv(cfg *Config) {
	cfg.PGSQL.Host = os.Getenv("PGSQL_HOST")
	cfg.PGSQL.User = os.Getenv("PGSQL_USER")
	cfg.PGSQL.Password = os.Getenv("PGSQL_PASS")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("pgsql", &cfg.PGSQL); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(folder string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}
