package metrics

import (
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/api/rules"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	logBasedMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "log_based_metric",
			Help: "Metrics based on logs",
		},
		[]string{"name", "metric", "namespace", "pod"},
	)

	activeMetrics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_log_based_metrics",
			Help: "Number of active log based metrics being tracked",
		},
		[]string{"placeholder"},
	)
)

func init() {
	prometheus.MustRegister(logBasedMetrics)
	prometheus.MustRegister(activeMetrics)
}

func IncActiveMetrics() {
	activeMetrics.With(prometheus.Labels{"placeholder": "placeholder"}).Inc()
}
func DecActiveMetrics() {
	activeMetrics.With(prometheus.Labels{"placeholder": "placeholder"}).Dec()
}

func IncrementMetric(rule rules.Rule, namespace, pod string) {
	logBasedMetrics.With(prometheus.Labels{"name": rule.Name, "metric": rule.Metric, "namespace": namespace, "pod": pod}).Inc()
}

func DecrementMetric(rule rules.Rule, namespace, pod string) {
	logBasedMetrics.With(prometheus.Labels{"name": rule.Name, "metric": rule.Metric, "namespace": namespace, "pod": pod}).Dec()
}

func SetMetric(rule rules.Rule, namespace, pod string, value float64) {
	logBasedMetrics.With(prometheus.Labels{"name": rule.Name, "metric": rule.Metric, "namespace": namespace, "pod": pod}).Set(value)
}

func ResetPrometheusGauges() {
	logBasedMetrics.Reset()
}
