package internal

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	maxRetries   = 5
	initialDelay = time.Second
	maxDelay     = time.Second * 5
)

func Backoff(fn func() error, skipErrors ...error) error {
	retries := 0
	delay := initialDelay

	for {
		err := fn()
		if err == nil {
			return nil
		}
		for _, skipError := range skipErrors {
			if errors.Is(err, skipError) {
				return err
			}
		}

		retries++
		if retries > maxRetries {
			log.Error().Err(err).Msg("maximum number of retries reached")
			return err
		}

		log.Error().Err(err).Msgf("backoff retrying in %s", delay)

		time.Sleep(delay)

		delay = delay + time.Second*2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
}
