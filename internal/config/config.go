package config

import (
	"github.com/afex/hystrix-go/hystrix"
	databaseInterface "github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/fajarardiyanto/flt-go-database/lib"
	env "github.com/fajarardiyanto/flt-go-env/lib"
	databaseListener "github.com/fajarardiyanto/flt-go-listener/interfaces"
	"github.com/fajarardiyanto/flt-go-logger/interfaces"
	log "github.com/fajarardiyanto/flt-go-logger/lib"
	jaeger "github.com/fajarardiyanto/flt-go-tracer/interfaces"
	"github.com/pkg/errors"
	"os"
)

var (
	config   Config
	logger   interfaces.Logger
	database databaseInterface.SQL
)

type Config struct {
	Timeout   int
	Namespace string
	Server    databaseListener.Server     `yaml:"server"`
	Database  databaseInterface.SQLConfig `yaml:"database"`
}

func GetLogger() interfaces.Logger {
	return logger
}

func GetConfig() *Config {
	return &config
}

func GetDBConn() databaseInterface.SQL {
	return database
}

func init() {
	if err := env.LoadEnv(".env"); err != nil {
		causer := errors.Cause(err)
		if os.IsNotExist(causer) {
			GetLogger().Info("Using default env config")
		} else {
			GetLogger().Error(causer).Quit()
		}
	}

	GetConfig().Namespace = "AFAIK Service News"

	GetConfig().Timeout = env.EnvInt("TIMEOUT", 30000)

	GetConfig().Database = databaseInterface.SQLConfig{
		Enable:        env.EnvBool("DATABASE_ENABLE", false),
		Driver:        env.EnvString("DATABASE_DRIVER", "mysql"),
		Host:          env.EnvString("DATABASE_HOST", "127.0.0.1"),
		Port:          env.EnvInt("DATABASE_PORT", 3306),
		Username:      env.EnvString("DATABASE_USERNAME", "user"),
		Password:      env.EnvString("DATABASE_PASSWORD", "user"),
		Database:      env.EnvString("DATABASE_NAME", "dbname"),
		AutoReconnect: env.EnvBool("DATABASE_AUTO_RECONNECT", true),
		StartInterval: env.EnvInt("DATABASE_INTERVAL", 2),
	}

	GetConfig().Server = databaseListener.Server{
		Name:    GetConfig().Namespace,
		Host:    env.EnvString("LISTENER_HOST", "127.0.0.1"),
		Port:    env.EnvInt("LISTENER_PORT", 8081),
		Timeout: env.EnvInt("LISTENER_TIMEOUT", 30000),
		Jaeger: jaeger.JaegerConfig{
			Host:   env.EnvString("JAEGER_HOST", "0.0.0.0"),
			Port:   env.EnvString("JAEGER_PORT", "6831"),
			Enable: env.EnvBool("JAEGER_ENABLE", false),
		},
	}
}

func Init() {
	logger = log.NewLib()
	logger.Init(GetConfig().Namespace)

	hystrix.ConfigureCommand("command_config", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  300,
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  50,
	})

	db := lib.NewLib()
	db.Init(logger)

	database = db.LoadSQLDatabase(GetConfig().Database)
	if err := database.LoadSQL(); err != nil {
		GetLogger().Error(err).Quit()
	}
}
