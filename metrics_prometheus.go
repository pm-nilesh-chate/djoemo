package djoemo

import (
	"context"
	"maps"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type prometheusmetrics struct {
	registry *prometheus.Registry
	// cfg           config.PrometheusMetrics
	queryCount    map[string]*prometheus.CounterVec
	queryDuration map[string]*prometheus.HistogramVec
}

func (m *prometheusmetrics) newCounter(caller string) *prometheus.CounterVec {
	opts := prometheus.CounterOpts{
		Name: strings.ToLower(caller),
		Help: "counter for function " + caller,
	}
	counter := prometheus.NewCounterVec(opts, []string{statusLabel})
	m.registry.MustRegister(counter)
	return counter
}

func (m *prometheusmetrics) newHistogramVec(caller string) *prometheus.HistogramVec {
	opts := prometheus.HistogramOpts{
		Name: strings.ToLower(caller),
		Help: "histogram duration for function " + caller,
		// WARNING: reduce the buckets after initial analysis
		Buckets: []float64{4, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 120},
	}
	// WARNING: add high cardinality labels like sdkhash, etc with caution
	histogram := prometheus.NewHistogramVec(opts, []string{statusLabel})
	m.registry.MustRegister(histogram)
	return histogram
}

const (
	statusLabel = "status"
	callerLabel = "caller" // NOTE: used separate metrics for now
	sourceLabel = "source"
	tableLabel  = "table"
)

func NewPrometheusMetrics(registry *prometheus.Registry) *prometheusmetrics {
	m := &prometheusmetrics{
		registry:      registry,
		queryCount:    make(map[string]*prometheus.CounterVec),
		queryDuration: make(map[string]*prometheus.HistogramVec),
	}
	return m
}

func (m *prometheusmetrics) Record(ctx context.Context, caller string, key KeyInterface, duration time.Duration, success bool) {
	if m.queryCount[caller] == nil || m.queryDuration[caller] == nil {
		m.queryCount[caller] = m.newCounter(caller)
		m.queryDuration[caller] = m.newHistogramVec(caller)
	}

	status := StatusFailure
	if success {
		status = StatusSuccess
	}

	labels := prometheus.Labels{statusLabel: status}
	if key.TableName() != "" {
		labels[tableLabel] = strings.ToLower(key.TableName())
	}

	maps.Copy(labels, LabelsFromContext(ctx))

	m.queryCount[caller].With(labels).Inc()
	m.queryDuration[caller].With(labels).Observe(float64(duration))
}
