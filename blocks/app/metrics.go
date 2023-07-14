package app

import (
	"github.com/blcvn/lib-golang-test/blocks/consensus/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// A Counter represents a monotonically increasing value.
type MetricCounter struct {
	// metrics.Counter
	prometheus.Counter
}

// With is used to provide label values when updating a Counter. This must be
// used to provide values for all LabelNames provided to CounterOpts.
func (s *MetricCounter) With(labelValues ...string) metrics.Counter {
	return s
}
func (s *MetricCounter) Add(delta float64) {
	s.Counter.Add(delta)
}

// A Gauge is a meter that expresses the current value of some metric.
type MetricGauge struct {
	prometheus.Gauge
}

// With is used to provide label values when updating a Counter. This must be
// used to provide values for all LabelNames provided to CounterOpts.
func (s *MetricGauge) With(labelValues ...string) metrics.Gauge {
	return s
}
func (s *MetricGauge) Add(delta float64) {
	s.Add(delta)
}

func (s *MetricGauge) Set(value float64) {
	s.Gauge.Set(value)
}

// A Histogram is a meter that records an observed value into quantized
// buckets.
type MetricHistogram struct {
	// metrics.Counter
	prometheus.Histogram
}

func (s *MetricHistogram) With(labelValues ...string) metrics.Histogram {
	return s
}
func (s *MetricHistogram) Observe(value float64) {
	s.Histogram.Observe(value)
}
