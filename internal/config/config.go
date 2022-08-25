package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	Config struct {
		HTTP  HTTPConfig
		PGSQL PGSQLConfig
		Auth  AuthConfig
		SMTP  SMTPConfig
		Redis RedisConfig
	}

	HTTPConfig struct {
		FrontendHost       string
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
		AllowedOrigins     []string
	}

	PGSQLConfig struct {
		Host         string
		User         string
		Password     string
		DatabaseName string `mapstructure:"dbname"`
		SSLMode      string `mapstructure:"sslmode"`
		Port         string
	}

	RedisConfig struct {
		Host      string
		Password  string
		DefaultDB int `mapstructure:"defaultDB"`
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

	SMTPConfig struct {
		Host     string
		User     string
		Password string
		Port     int
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

	if err := setEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setEnv(cfg *Config) error {
	var ok bool
	// HTTP
	cfg.HTTP.Host, ok = os.LookupEnv("HTTP_HOST")
	if !ok {
		return errors.New("empty http host env")
	}
	cfg.HTTP.FrontendHost, ok = os.LookupEnv("HTTP_FRONTEND_HOST")
	if !ok {
		return errors.New("empty frontend host env")
	}

	val, ok := os.LookupEnv("ALLOWED_ORIGINS")
	if !ok || len(val) < 1 {
		return errors.New("wrong allowed origins env")
	}
	origins, err := parseCORS(val)
	if err != nil {
		return err
	}
	cfg.HTTP.AllowedOrigins = origins

	// PostgresSQL connection
	cfg.PGSQL.Host, ok = os.LookupEnv("POSTGRES_HOST")
	if !ok {
		return errors.New("empty postgres host env")
	}
	cfg.PGSQL.User, ok = os.LookupEnv("POSTGRES_USER")
	if !ok {
		return errors.New("empty postgres user env")
	}
	cfg.PGSQL.Password, ok = os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		return errors.New("empty postgres password env")
	}
	val, ok = os.LookupEnv("POSTGRES_PORT")
	if !ok {
		cfg.PGSQL.Port = "5432"
	} else {
		cfg.PGSQL.Port = val
	}

	// JWT
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	// Password salt
	cfg.Auth.PasswordSalt = os.Getenv("PASS_SALT")

	// SMTP
	cfg.SMTP.Host, ok = os.LookupEnv("SMTP_HOST")
	if !ok {
		return errors.New("empty smtp host env")
	}
	cfg.SMTP.User, ok = os.LookupEnv("SMTP_USER")
	if !ok {
		return errors.New("empty smtp user env")
	}
	cfg.SMTP.Password, ok = os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		return errors.New("empty smtp password env")
	}
	val, ok = os.LookupEnv("SMTP_PORT")
	if !ok {
		return errors.New("empty smtp port env")
	}
	port, err := strconv.Atoi(val)
	if err != nil {
		return err
	} else {
		cfg.SMTP.Port = port
	}

	// Redis
	cfg.Redis.Host, ok = os.LookupEnv("REDIS_HOST")
	if !ok {
		return errors.New("empty redis host env")
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

	return nil
}

func parseCORS(env string) ([]string, error) {
	out := make([]string, 0)
	arr := strings.Split(env, ",")
	if len(arr) < 1 {
		return nil, errors.New("wrong allowed origins env")
	}

	for i, item := range arr {
		pattern := `^https?:\/\/\w+(\.\w+)*(:[0-9]+)?(\/.*)?$`
		if matched, _ := regexp.MatchString(pattern, item); !matched {
			return nil, fmt.Errorf("wrong allowed origins number %d", i+1)
		}
		out = append(out, item)
	}

	return out, nil
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

	if err := viper.UnmarshalKey("redis", &cfg.Redis); err != nil {
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
