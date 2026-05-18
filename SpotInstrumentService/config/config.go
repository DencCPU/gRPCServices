package spotconfig

import "time"

type Config struct {
	Redis         Redis         `mapstructure:"redis"`
	Server        Server        `mapstructure:"server"`
	Storage       Storage       `mapstructure:"storage"`
	OtelCollector OtelCollector `mapstructure:"collector"`
}

type Redis struct {
	Host         string        `mapstructure:"host"`
	Port         string        `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_id_le_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	Timeout      time.Duration `mapstructure:"timeout"`
}

type Server struct {
	Port                  string `mapstructure:"port"`
	Host                  string `mapstructure:"host"`
	Network               string `mapstructure:"network"`
	RequestPerSecondLimit uint   `mapstructure:"request_per_second_limit"`
}

type Storage struct {
	Timeout time.Duration `mapstructure:"work_time"`
}

type OtelCollector struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	TracePercentage int           `mapstructure:"trace_percentage"`
	MetricInterval  time.Duration `mapstructure:"metric_interval"`
}
