package app

import (
	"log"
	"log/slog"
	"os"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Config Конфиг приложения
type Config struct {
	ProductServiceURL        string `validate:"required" env:"PRODUCT_SERVICE_URL"`
	ProductServiceToken      string `validate:"required" env:"PRODUCT_SERVICE_TOKEN"`
	LomsServiceAddr          string `validate:"required" env:"LOMS_SERVICE_ADDR"`
	ProductServiceRetryCount int    `env:"PRODUCT_SERVICE_RETRY_COUNT"`
	ProductServiceRPSLimit   int64  `env:"PRODUCT_SERVICE_RPS_LIMIT"`
	ListenAddr               string `validate:"required" env:"LISTEN_ADDR"`
	Logger                   *slog.Logger
	TracerProvider           trace.TracerProvider
	MeterProvider            metric.MeterProvider
}

// ReadConfig Конфигурация приложения (мини-viper :))
func ReadConfig() *Config {
	conf := &Config{}
	v := reflect.ValueOf(conf).Elem()
	for i := 0; i < v.NumField(); i++ {
		envVar, exists := v.Type().Field(i).Tag.Lookup("env")
		if !exists {
			continue
		}
		envVal, exists := os.LookupEnv(envVar)
		if !exists {
			log.Fatalf("env not found: %s", envVar)
		}
		switch v.Field(i).Interface().(type) {
		case string:
			v.Field(i).SetString(envVal)
		case int:
		case int64:
			intVal, err := strconv.Atoi(envVal)
			if err != nil {
				log.Fatalf("invalid value for env: %s value: %s", envVar, envVal)
			}
			v.Field(i).SetInt(int64(intVal))
		default:
			log.Fatalf("unknown type for env: %s type: %s", envVar, v.Field(i).Type().Name())
		}
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(conf); err != nil {
		log.Fatal(err)
	}

	return conf
}
