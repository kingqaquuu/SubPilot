package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	API      APIConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env  string
	Name string
	Port string
}

type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type APIConfig struct {
	Prefix string
}

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setDefaults(v)
	if err := bindEnv(v); err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{
		App: AppConfig{
			Env:  v.GetString("app.env"),
			Name: v.GetString("app.name"),
			Port: v.GetString("app.port"),
		},
		Server: ServerConfig{
			ReadTimeout:  v.GetDuration("server.read_timeout"),
			WriteTimeout: v.GetDuration("server.write_timeout"),
		},
		API: APIConfig{
			Prefix: v.GetString("api.prefix"),
		},
		Postgres: PostgresConfig{
			Host:     v.GetString("postgres.host"),
			Port:     v.GetString("postgres.port"),
			Database: v.GetString("postgres.database"),
			User:     v.GetString("postgres.user"),
			Password: v.GetString("postgres.password"),
			SSLMode:  v.GetString("postgres.ssl_mode"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("redis.host"),
			Port:     v.GetString("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		JWT: JWTConfig{
			Secret:    v.GetString("jwt.secret"),
			ExpiresIn: v.GetDuration("jwt.expires_in"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.env", "development")
	v.SetDefault("app.name", "SubPilot")
	v.SetDefault("app.port", "18080")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")
	v.SetDefault("api.prefix", "/api/v1")
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.database", "subpilot")
	v.SetDefault("postgres.user", "subpilot")
	v.SetDefault("postgres.password", "")
	v.SetDefault("postgres.ssl_mode", "disable")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.secret", "")
	v.SetDefault("jwt.expires_in", "24h")
}

func bindEnv(v *viper.Viper) error {
	envBindings := map[string]string{
		"app.env":              "APP_ENV",
		"app.name":             "APP_NAME",
		"app.port":             "APP_PORT",
		"server.read_timeout":  "SERVER_READ_TIMEOUT",
		"server.write_timeout": "SERVER_WRITE_TIMEOUT",
		"api.prefix":           "API_PREFIX",
		"postgres.host":        "POSTGRES_HOST",
		"postgres.port":        "POSTGRES_PORT",
		"postgres.database":    "POSTGRES_DB",
		"postgres.user":        "POSTGRES_USER",
		"postgres.password":    "POSTGRES_PASSWORD",
		"postgres.ssl_mode":    "POSTGRES_SSL_MODE",
		"redis.host":           "REDIS_HOST",
		"redis.port":           "REDIS_PORT",
		"redis.password":       "REDIS_PASSWORD",
		"redis.db":             "REDIS_DB",
		"jwt.secret":           "JWT_SECRET",
		"jwt.expires_in":       "JWT_EXPIRES_IN",
	}

	for key, env := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			return fmt.Errorf("bind env %s: %w", env, err)
		}
	}

	return nil
}

func (c *Config) Validate() error {
	if c.JWT.Secret == "change-me-to-a-long-random-secret" {
		return errors.New("JWT_SECRET must be set to a non-default value")
	}

	if c.App.Env != "production" {
		return nil
	}

	if c.Postgres.Password == "" || c.Postgres.Password == "change-me" {
		return errors.New("POSTGRES_PASSWORD must be set to a non-default value in production")
	}
	if c.JWT.Secret == "" {
		return errors.New("JWT_SECRET must be set to a non-default value in production")
	}

	return nil
}
