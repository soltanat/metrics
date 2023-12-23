package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		log = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	})

	return log
}
