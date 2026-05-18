package orderconfig

import "time"

type Config struct {
	Server         Server         `mapstructure:"server"`
	Postgres       Postgres       `mapstructure:"postgres"`
	Notify         Notify         `mapstructure:"notify"`
	BreakerSetting BreakerSetting `mapstructure:"breakersetting"`
	OtelCollector  OtelCollector  `mapstructure:"collector"`
}

type Server struct {
	Port                  int    `mapstructure:"port"`
	Host                  string `mapstructure:"host"`
	Network               string `mapstructure:"network"`
	RequestPerSecondLimit uint   `mapstructure:"request_per_second_limit"`
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
