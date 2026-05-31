package orderconfig

import "time"

type Config struct {
	Server         Server         `mapstructure:"server"`
	Postgres       Postgres       `mapstructure:"postgres"`
	Notify         Notify         `mapstructure:"notify"`
	BreakerSetting BreakerSetting `mapstructure:"breakersetting"`
	OtelCollector  OtelCollector  `mapstructure:"collector"`
	Kafka          Kafka          `mapstructure:"kafka"`
}

type Server struct {
	Port                    int           `mapstructure:"port"`
	Host                    string        `mapstructure:"host"`
	Network                 string        `mapstructure:"network"`
	RequestPerSecondLimit   uint          `mapstructure:"request_per_second_limit"`
	ClientConnectionTimeout time.Duration `mapstructure:"client_connection_timeout"`
}

type Postgres struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	User                string        `mapstructure:"user"`
	Password            string        `mapstructure:"password"`
	Name                string        `mapstructure:"name"`
	Sslmode             string        `mapstructure:"sslmode"`
	ControlChanSize     int           `mapstructure:"chan_size"`
	IdempotencyCacheTTL time.Duration `mapstructure:"idepmpotency_cache_ttl"`
	MarketCacheTTL      time.Duration `mapstructure:"market_cache_ttl"`
}

type Notify struct {
	TickerInterval time.Duration `mapstructure:"ticker_interval"`
}

type BreakerSetting struct {
	Name           string        `mapstructure:"name"`
	MaxRequests    uint32        `mapstructure:"max_request"`
	Interval       time.Duration `mapstructure:"interval"`
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxFailRequest uint32        `mapstructure:"max_fail_request"`
}

type OtelCollector struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	TracePercentage int           `mapstructure:"trace_percentage"`
	MetricInterval  time.Duration `mapstructure:"metric_interval"`
}

type Kafka struct {
	Brokers  []string      `mapstructure:"brokers"`
	Topic    string        `mapstructure:"topic"`
	GroupID  string        `mapstructure:"group_id"`
	MinBytes int           `mapstructure:"min_bytes"`
	MaxBytes int           `mapstructure:"max_bytes"`
	MaxWait  time.Duration `mapstructure:"max_wait"`
}
