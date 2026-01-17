package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig        `mapstructure:"app"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	SMTP       SMTPConfig       `mapstructure:"smtp"`
	Slack      SlackConfig      `mapstructure:"slack"`
}

type AppConfig struct {
	Name     string `mapstructure:"name"`
	Env      string `mapstructure:"env"`
	Port     int    `mapstructure:"port"`
	LogLevel int    `mapstructure:"log_level"`
}

type SlackConfig struct {
	WebhookURL string `mapstructure:"webhook_url"`
}

type SMTPConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	From string `mapstructure:"from"`
}

type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Params          string        `mapstructure:"params"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
		m.Params,
	)
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type KafkaConfig struct {
	Brokers []string          `mapstructure:"brokers"`
	Topics  map[string]string `mapstructure:"topics"`
}

type PrometheusConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	// env override support (APP_PORT, MYSQL_HOST, etc.)
	v.SetEnvPrefix("NOTIF")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
