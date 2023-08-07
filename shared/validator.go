package shared

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var once sync.Once
var v *validator.Validate

// GetValidator is responsible for returning a single instance of the validator.
func GetValidator() *validator.Validate {
	once.Do(func() {
		log.Info().Msg("Validator initialized.")
		v = validator.New()
	})

	return v
}
