package config

import (
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

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return &Config{
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
	}, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.env", "development")
	v.SetDefault("app.name", "SubPilot")
	v.SetDefault("app.port", "8080")
	v.SetDefault("server.read_timeout", "10s")
	v.SetDefault("server.write_timeout", "10s")
	v.SetDefault("api.prefix", "/api/v1")
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.database", "subpilot")
	v.SetDefault("postgres.user", "subpilot")
	v.SetDefault("postgres.password", "")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.secret", "")
	v.SetDefault("jwt.expires_in", "24h")
}
