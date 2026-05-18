package apiconfig

import "time"

type Config struct {
	Server         Server         `mapstructure:"server"`
	BreakerSetting BreakerSetting `mapstructure:"breakersetting"`
	OtelCollector  OtelCollector  `mapstructure:"collector"`
}

type Server struct {
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
	Network string `mapstructure:"network"`
}

type BreakerSetting struct {
	Name           string        `mapstructure:"name"`             //Название брейкера
	MaxRequests    uint32        `mapstructure:"max_request"`      //Максимально кол-во запросов, пропускаемых в полуоткрытом режиме(Half-open)
	Interval       time.Duration `mapstructure:"interval"`         //Период сброса статистики подсчета неудачных запросов в закрытом режиме(Close). Время в секундах.
	Timeout        time.Duration `mapstructure:"timeout"`          //Время прибывания брейкера в открытом состоянии, перед переходов в Half-open.
	MaxFailRequest uint32        `mapstructure:"max_fail_request"` //Количество неудачных запросов, после которого брейкер перейдет в состояние Open
}

type OtelCollector struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	TracePercentage int           `mapstructure:"trace_percentage"`
	MetricInterval  time.Duration `mapstructure:"metric_interval"`
}
