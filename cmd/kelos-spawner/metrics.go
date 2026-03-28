package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	// discoveryTotal counts the total number of completed discovery cycles.
	discoveryTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kelos_spawner_discovery_total",
			Help: "Total number of completed discovery cycles",
		},
	)

	// discoveryErrorsTotal counts the total number of failed discovery cycles.
	discoveryErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kelos_spawner_discovery_errors_total",
			Help: "Total number of failed discovery cycles",
		},
	)

	// itemsDiscoveredTotal counts the total number of work items discovered.
	itemsDiscoveredTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kelos_spawner_items_discovered_total",
			Help: "Total number of work items discovered",
		},
	)

	// tasksCreatedTotal counts the total number of Tasks created by the spawner.
	tasksCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kelos_spawner_tasks_created_total",
			Help: "Total number of Tasks created by the spawner",
		},
	)

	// discoveryDurationSeconds records the duration of discovery cycles.
	discoveryDurationSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "kelos_spawner_discovery_duration_seconds",
			Help:    "Duration of discovery cycles in seconds",
			Buckets: []float64{0.5, 1, 2, 5, 10, 30, 60},
		},
	)
)

func init() {
	metrics.Registry.MustRegister(
		discoveryTotal,
		discoveryErrorsTotal,
		itemsDiscoveredTotal,
		tasksCreatedTotal,
		discoveryDurationSeconds,
	)
}
