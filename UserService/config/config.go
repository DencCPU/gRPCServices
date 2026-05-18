package userconfig

import "time"

type Config struct {
	Server        Server        `mapstructure:"server"`
	Postgres      Postgres      `mapstructure:"postgres"`
	JWT           JWT           `mapstructure:"jwt"`
	OtelCollector OtelCollector `mapstructure:"collector"`
}

type Server struct {
	Host                  string `mapstructure:"host"`
	Port                  string `mapstructure:"port"`
	Network               string `mapstructure:"network"`
	RequestPerSecondLimit uint   `mapstructure:"request_per_second_limit"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Sslmode  string `mapstructure:"sslmode"`
}

type JWT struct {
	Secret string `mapstructure:"secret"`
	TTL    int    `mapstructure:"ttl"`
}

type OtelCollector struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	TracePercentage int           `mapstructure:"trace_percentage"`
	MetricInterval  time.Duration `mapstructure:"metric_interval"`
}
