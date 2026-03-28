package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
)

func TestMetricsRegistered(t *testing.T) {
	collectors := []struct {
		name      string
		collector prometheus.Collector
	}{
		{"discoveryTotal", discoveryTotal},
		{"discoveryErrorsTotal", discoveryErrorsTotal},
		{"itemsDiscoveredTotal", itemsDiscoveredTotal},
		{"tasksCreatedTotal", tasksCreatedTotal},
		{"discoveryDurationSeconds", discoveryDurationSeconds},
	}

	for _, tc := range collectors {
		t.Run(tc.name, func(t *testing.T) {
			ch := make(chan *prometheus.Desc, 10)
			tc.collector.Describe(ch)
			close(ch)
			if len(ch) == 0 {
				t.Errorf("expected at least one descriptor for %s", tc.name)
			}
		})
	}
}

func TestDiscoveryTotalCounter(t *testing.T) {
	before := testutil.ToFloat64(discoveryTotal)
	discoveryTotal.Inc()
	after := testutil.ToFloat64(discoveryTotal)
	if after != before+1 {
		t.Errorf("expected discoveryTotal to increment by 1, got delta %f", after-before)
	}
}

func TestDiscoveryErrorsTotalCounter(t *testing.T) {
	before := testutil.ToFloat64(discoveryErrorsTotal)
	discoveryErrorsTotal.Inc()
	after := testutil.ToFloat64(discoveryErrorsTotal)
	if after != before+1 {
		t.Errorf("expected discoveryErrorsTotal to increment by 1, got delta %f", after-before)
	}
}

func TestItemsDiscoveredTotalCounter(t *testing.T) {
	before := testutil.ToFloat64(itemsDiscoveredTotal)
	itemsDiscoveredTotal.Add(5)
	after := testutil.ToFloat64(itemsDiscoveredTotal)
	if after != before+5 {
		t.Errorf("expected itemsDiscoveredTotal to increment by 5, got delta %f", after-before)
	}
}

func TestTasksCreatedTotalCounter(t *testing.T) {
	before := testutil.ToFloat64(tasksCreatedTotal)
	tasksCreatedTotal.Add(3)
	after := testutil.ToFloat64(tasksCreatedTotal)
	if after != before+3 {
		t.Errorf("expected tasksCreatedTotal to increment by 3, got delta %f", after-before)
	}
}

func TestDiscoveryDurationSecondsHistogram(t *testing.T) {
	discoveryDurationSeconds.Observe(1.5)

	ch := make(chan prometheus.Metric, 1)
	discoveryDurationSeconds.Collect(ch)
	m := <-ch

	var dto dto.Metric
	if err := m.Write(&dto); err != nil {
		t.Fatalf("writing metric: %v", err)
	}
	h := dto.GetHistogram()
	if h.GetSampleCount() == 0 {
		t.Error("expected sample count > 0 after Observe")
	}
	if h.GetSampleSum() < 1.5 {
		t.Errorf("expected sample sum >= 1.5, got %f", h.GetSampleSum())
	}
}
