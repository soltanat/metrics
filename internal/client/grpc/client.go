package grpc

import (
	"context"
	"github.com/soltanat/metrics/internal/model"
	"github.com/soltanat/metrics/internal/proto"
	"google.golang.org/grpc"
)

type Client struct {
	cli proto.MetricsServiceClient
}

func NewClient(conn grpc.ClientConnInterface) *Client {
	cli := proto.NewMetricsServiceClient(conn)
	return &Client{
		cli: cli,
	}
}

func (c *Client) Send(ctx context.Context, m *model.Metric) error {
	//TODO implement me
	panic("implement me")
}

func (c *Client) Updates(ctx context.Context, mm []model.Metric) error {
	req := &proto.BulkSendMetricRequest{
		Metrics: make([]*proto.SendMetricRequest, 0, len(mm)),
	}

	for i := 0; i < len(mm); i++ {
		m := &mm[i]
		req.Metrics = append(req.Metrics, &proto.SendMetricRequest{
			Id:    m.Name,
			Type:  proto.MetricType(m.Type),
			Delta: &m.Counter,
			Value: &m.Gauge,
		})
	}

	_, err := c.cli.BulkSendMetric(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Update(ctx context.Context, m *model.Metric) error {
	//TODO implement me
	panic("implement me")
}
