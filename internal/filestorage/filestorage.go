package filestorage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/storage"
)

type FileStorage struct {
	storage.Storage
	store    func(m *model.Metric) error
	file     *os.File
	mu       *sync.Mutex
	interval time.Duration
	stopCh   chan struct{}
	closeCh  chan struct{}
}

func New(storage storage.Storage, interval time.Duration, path string, restore bool) (*FileStorage, error) {
	s := &FileStorage{
		Storage:  storage,
		mu:       &sync.Mutex{},
		interval: interval,
		stopCh:   make(chan struct{}),
		closeCh:  make(chan struct{}),
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	s.file = file

	if restore {
		err = s.restore()
		if err != nil {
			return nil, err
		}
	} else {
		err = s.file.Truncate(0)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *FileStorage) Store(m *model.Metric) error {
	err := s.Storage.Store(m)
	if err != nil {
		return fmt.Errorf("failed to store: %w", err)
	}
	if s.interval == 0 {
		return s.flush()
	}
	return nil
}

func (s *FileStorage) restore() error {
	l := logger.Get()

	l.Info().Msg("file storage restored started")
	start := time.Now()
	defer func() {
		l.Info().Dur("duration", time.Since(start)).Msg("file storage restored")
	}()

	_, err := s.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	dec := json.NewDecoder(s.file)

	for dec.More() {
		var m model.Metric
		err := dec.Decode(&m)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}
		err = s.Storage.Store(&m)
		if err != nil {
			return fmt.Errorf("failed to store: %w", err)
		}
	}

	return nil
}

func (s *FileStorage) flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ms, err := s.Storage.GetList()
	if err != nil {
		return fmt.Errorf("failed to get list: %w", err)
	}

	err = s.file.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate: %w", err)
	}
	_, err = s.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}

	enc := json.NewEncoder(s.file)
	for _, m := range ms {
		err := enc.Encode(m)
		if err != nil {
			return fmt.Errorf("failed to encode: %w", err)
		}
	}

	return nil
}

func (s *FileStorage) Stop() error {
	s.stopCh <- struct{}{}
	<-s.closeCh
	return s.file.Close()
}

func (s *FileStorage) Start() error {
	if s.interval == 0 {
		return fmt.Errorf("interval is zero")
	}
	go func() {
		l := logger.Get()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.flush(); err != nil {
					l.Error().Err(err).Msg("file storage error")
				} else {
					l.Info().Msg("file storage flushed")
				}
			case <-s.stopCh:
				if err := s.flush(); err != nil {
					l.Error().Err(err).Msg("file storage error")
				} else {
					l.Info().Msg("file storage flushed")
				}
				s.closeCh <- struct{}{}
				l.Info().Msg("file storage stopped")
				return
			}
		}

	}()
	return nil
}
