package metrics

import (
	"github.com/hyperledger/fabric/common/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Add increments a counter value.

type AppMetricProvider struct {
	metrics.Provider
}

// NewCounter creates a new instance of a Counter.
func (s *AppMetricProvider) NewCounter(opt metrics.CounterOpts) metrics.Counter {
	pCounterOpts := prometheus.CounterOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: map[string]string{},
	}
	mt := prometheus.NewCounter(pCounterOpts)
	return &MetricCounter{
		Counter: mt,
	}
}

// NewGauge creates a new instance of a Gauge.
func (s *AppMetricProvider) NewGauge(opt metrics.GaugeOpts) metrics.Gauge {
	pGaugeOpts := prometheus.GaugeOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: make(map[string]string, 0),
	}
	mt := prometheus.NewGauge(pGaugeOpts)
	return &MetricGauge{
		Gauge: mt,
	}
}

// NewHistogram creates a new instance of a Histogram.
func (s *AppMetricProvider) NewHistogram(opt metrics.HistogramOpts) metrics.Histogram {
	pHistogramOpts := prometheus.HistogramOpts{
		Namespace:                       opt.Namespace,
		Subsystem:                       opt.Subsystem,
		Name:                            opt.Name,
		Help:                            opt.Help,
		ConstLabels:                     map[string]string{},
		Buckets:                         opt.Buckets,
		NativeHistogramBucketFactor:     0,
		NativeHistogramZeroThreshold:    0,
		NativeHistogramMaxBucketNumber:  0,
		NativeHistogramMinResetDuration: 0,
		NativeHistogramMaxZeroThreshold: 0,
	}
	mt := prometheus.NewHistogram(pHistogramOpts)
	return &MetricHistogram{
		Histogram: mt,
	}
}
