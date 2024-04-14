package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	HTTP  HTTP  `mapstructure:"http"`
	CACHE CACHE `mapstructure:"cache"`
	PG    PG    `mapstructure:"postgres"`
}

type HTTP struct {
	Port           string        `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	IdleTimeout    time.Duration `mapstructure:"idle_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

type CACHE struct {
	TTL         time.Duration `mapstructure:"ttl"`
	MaxCost     int64         `mapstructure:"max_cost"`
	NumCounters int64         `mapstructure:"num_counters"`
}

type PG struct {
	URL             string        `mapstructure:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func NewConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath("./config")
	v.AutomaticEnv()

	path, exists := os.LookupEnv("CONFIG_PATH")
	if exists {
		v.SetConfigFile(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Printf("Loaded configuration: %+v\n", cfg)

	return &cfg, nil
}

func RequiredConfig() *Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}
