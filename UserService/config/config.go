package userconfig

type Config struct {
	Server   Server   `mapstructure:"server"`
	Postgres Postgres `mapstructure:"postgres"`
	JWT      JWT      `mapstructure:"jwt"`
	Jaeger   Jaeger   `mapstructure:"jaeger"`
}

type Server struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Network string `mapstructure:"network"`
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

type Jaeger struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
