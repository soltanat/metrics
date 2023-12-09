package internal

type Poll interface {
	Get() []Metric
	Poll() error
}

type Reporter interface {
	Report() error
}
