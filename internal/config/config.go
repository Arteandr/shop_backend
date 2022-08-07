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
		Auth  AuthConfig
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
		Port         string
	}

	AuthConfig struct {
		PasswordSalt    string
		JWT             JWTConfig
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
	}

	JWTConfig struct {
		SigningKey string
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
	// PostgresSQL connection
	cfg.PGSQL.Host = os.Getenv("POSTGRES_HOST")
	cfg.PGSQL.User = os.Getenv("POSTGRES_USER")
	cfg.PGSQL.Password = os.Getenv("POSTGRES_PASSWORD")
	val, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		cfg.PGSQL.Port = "5432"
	} else {
		cfg.PGSQL.Port = val
	}

	// JWT
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	// Password salt
	cfg.Auth.PasswordSalt = os.Getenv("PASS_SALT")
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("pgsql", &cfg.PGSQL); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
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
