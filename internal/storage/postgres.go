package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/soltanat/metrics/internal/model"
)

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(conn *pgx.Conn) *PostgresStorage {
	return &PostgresStorage{conn: conn}
}

func (s *PostgresStorage) Store(metric *model.Metric) error {
	if metric.Type == model.MetricTypeCounter {
		_, err := s.conn.Exec(
			context.Background(),
			`INSERT INTO metrics.metrics_counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2`, metric.Name, metric.Counter,
		)
		if err != nil {
			return err
		}
	} else if metric.Type == model.MetricTypeGauge {
		_, err := s.conn.Exec(
			context.Background(),
			"INSERT INTO metrics.metrics_gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", metric.Name, metric.Gauge,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresStorage) StoreBatch(metrics []model.Metric) error {
	ctx := context.Background()
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, m := range metrics {
		if m.Type == model.MetricTypeGauge {
			batch.Queue(`INSERT INTO metrics.metrics_gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2`, m.Name, m.Gauge)
		} else if m.Type == model.MetricTypeCounter {
			batch.Queue(`INSERT INTO metrics.metrics_counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2`, m.Name, m.Counter)
		}
	}

	br := tx.SendBatch(ctx, batch)

	err = br.Close()
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStorage) GetGauge(name string) (*model.Metric, error) {
	row := s.conn.QueryRow(context.Background(), "SELECT value FROM metrics.metrics_gauge WHERE name = $1", name)
	var v float64
	err := row.Scan(&v)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrMetricNotFound
		}
		return nil, err
	}
	return model.NewGauge(name, v), nil
}

func (s *PostgresStorage) GetCounter(name string) (*model.Metric, error) {
	row := s.conn.QueryRow(context.Background(), "SELECT value FROM metrics.metrics_counter WHERE name = $1", name)
	var v int64
	err := row.Scan(&v)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrMetricNotFound
		}
		return nil, err
	}
	return model.NewCounter(name, v), nil
}

func (s *PostgresStorage) GetList() ([]model.Metric, error) {
	row, err := s.conn.Query(context.Background(), "SELECT name, value FROM metrics.metrics_gauge")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	metrics := make([]model.Metric, 0)
	for row.Next() {
		var name string
		var v float64
		err = row.Scan(&name, &v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, *model.NewGauge(name, v))
	}

	row, err = s.conn.Query(context.Background(), "SELECT name, value FROM metrics.metrics_counter")
	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var name string
		var v int64
		err = row.Scan(&name, &v)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, *model.NewCounter(name, v))
	}

	return metrics, nil
}
