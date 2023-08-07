package configs

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config is a struct that will receive configuration options via environment
// variables.
type Config struct {
	App struct {
		CORS struct {
			AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"`
			AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS"`
			AllowedMethods   []string `mapstructure:"ALLOWED_METHODS"`
			AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS"`
			Enable           bool     `mapstructure:"ENABLE"`
			MaxAgeSeconds    int      `mapstructure:"MAX_AGE_SECONDS"`
		}
		Name     string `mapstructure:"NAME"`
		Revision string `mapstructure:"REVISION"`
		URL      string `mapstructure:"URL"`
	}

	Cache struct {
		Redis struct {
			Primary struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Password string `mapstructure:"PASSWORD"`
			}
		}
	}

	DB struct {
		MySQL struct {
			Read struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Username string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			}
			Write struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Username string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			}
		}
	}

	Event struct {
		Consumer struct {
			SQS struct {
				AccessKeyID       string `mapstructure:"ACCESS_KEY_ID"`
				BackoffSeconds    int    `mapstructure:"BACKOFF_SECONDS"`
				MaxMessage        int64  `mapstructure:"MAX_MESSAGE"`
				MaxRetries        int    `mapstructure:"MAX_RETRIES"`
				MaxRetriesConsume int    `mapstructure:"MAX_RETRIES_CONSUME"`
				Region            string `mapstructure:"REGION"`
				SecretAccessKey   string `mapstructure:"SECRET_ACCESS_KEY"`
				WaitTimeSeconds   int64  `mapstructure:"WAIT_TIME_SECONDS"`

				Topics struct {
					FooBarBaz struct {
						Enabled bool   `mapstructure:"ENABLED"`
						URL     string `mapstructure:"URL"`
					} `mapstructure:"FOOBARBAZ"`
				}
			}
		}

		Producer struct {
			SNS struct {
				AccessKeyID     string `mapstructure:"ACCESS_KEY_ID"`
				MaxRetries      int    `mapstructure:"MAX_RETRIES"`
				Region          string `mapstructure:"REGION"`
				SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
				Topics          struct {
					FooCreated struct {
						ARN     string `mapstructure:"ARN"`
						Enabled bool   `mapstructure:"ENABLED"`
					} `mapstructure:"FOO_CREATED"`
				}
			}
		}
	}

	Server struct {
		Env      string `mapstructure:"ENV"`
		LogLevel string `mapstructure:"LOG_LEVEL"`
		Port     string `mapstructure:"PORT"`
		Shutdown struct {
			CleanupPeriodSeconds int64 `mapstructure:"CLEANUP_PERIOD_SECONDS"`
			GracePeriodSeconds   int64 `mapstructure:"GRACE_PERIOD_SECONDS"`
		}
	}
}

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading config file")
	}

	once.Do(func() {
		log.Info().Msg("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	})

	return &conf
}
