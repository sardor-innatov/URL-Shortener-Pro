package config

import (
	"fmt"
	"sync"
)

var envProject EnvProject


func init() {
	sync.OnceFunc(func() {
		MustLoad(&envProject, ".env")
	})()
}

func ProjectEnv() EnvProject {
	return envProject
}

type EnvProject struct {
	BaseURL    string `env:"BASE_URL"`

	PgDB       string `env:"POSTGRES_DB"`
	PgHost     string `env:"POSTGRES_HOST"`
	PgPort     uint   `env:"POSTGRES_PORT"`
	PgUser     string `env:"POSTGRES_USER"`
	PgPassword string `env:"POSTGRES_PASSWORD"`

	RedisAddr string `env:"REDIS_ADDR"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	JwtSecret  string `env:"JWT_SECRET"`
	JwtExpire  int64  `env:"JWT_EXPIRE"`
	JwtRefresh int64  `env:"JWT_REFRESH"`

	ClickWorkers int `env:"CLICK_WORKERS_COUNT"`
}

type Config struct {
	host     string
	user     string
	password string
	dbname   string
	port     uint
	sslmode  bool
}

func (c Config) ssl() string {
	if !c.sslmode {
		return "sslmode=disable"
	}
	return ""
}

func (c Config) build() string {
	return fmt.Sprintf(
		`host=%s
    	 port=%d 
    	 user=%s 
    	 password=%s 
    	 dbname=%s 
    	 %s`,
			c.host,
			c.port,
			c.user,
			c.password,
			c.dbname,
			c.ssl(),
	)
}

func NewConfig(
	host, user, password, dbname string, port uint,
) Config {
	return Config{
		host, user, password, dbname, port, false,
	}
}