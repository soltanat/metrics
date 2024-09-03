package grpchandler

import (
	"context"
	"errors"
	"github.com/soltanat/metrics/internal/logger"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/proto"
	"github.com/soltanat/metrics/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	proto.MetricsServiceServer
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) BulkSendMetric(ctx context.Context, request *proto.BulkSendMetricRequest) (*proto.BulkSendMetricResponse, error) {
	l := logger.Get()

	mm, err := model.NewMetrics(ConvertToInterfaces(request.Metrics))
	if err != nil {
		if errors.As(err, &model.ErrBadRequest{}) {
			l.Error().Err(err).Msg("Error parsing metrics")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	if err := s.storage.StoreBatch(mm); err != nil {
		l.Error().Msgf("Error storing metric: %s", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	l.Info().Msg("Metrics stored successfully")

	return &proto.BulkSendMetricResponse{}, nil
}

type InputMetric struct {
	m *proto.SendMetricRequest
}

func (i *InputMetric) Type() string {
	return i.m.Type.String()
}

func (i *InputMetric) Value() *float64 {
	return i.m.Value
}

func (i *InputMetric) Delta() *int64 {
	return i.m.Delta
}

func (i *InputMetric) ID() string {
	return i.m.Id
}

func ConvertToInterfaces(input []*proto.SendMetricRequest) []model.InputMetric {
	metrics := make([]model.InputMetric, len(input))
	for i, m := range input {
		metrics[i] = &InputMetric{m}
	}
	return metrics
}
