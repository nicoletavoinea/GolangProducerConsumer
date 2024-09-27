package functions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var myGeneralMetrics PrometheusMetricsGeneral
var myTypesMetrics PrometheusMetricsTypes

type PrometheusMetricsGeneral struct {
	ProcessedTasks prometheus.Gauge
	DoneTasks      prometheus.Counter
}

type PrometheusMetricsTypes struct {
	TotalTasks  *prometheus.CounterVec
	TotalValues *prometheus.CounterVec
}

func CreatePrometheusMetricsGeneral() {
	processed := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "processed_tasks",
			Help: "Total number of processed tasks.",
		},
	)

	done := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "done_tasks",
			Help: "Total number of done tasks.",
		},
	)

	// Register the metrics with Prometheus
	prometheus.MustRegister(processed)
	prometheus.MustRegister(done)

	myGeneralMetrics = PrometheusMetricsGeneral{
		ProcessedTasks: processed,
		DoneTasks:      done,
	}
}

func CreatePrometheusMetricsTypes() {
	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_tasks",
			Help: "Total number of tasks grouped by type",
		},
		[]string{"type"},
	)

	values := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_values",
			Help: "Sum of the values of the tasks grouped by type",
		},
		[]string{"type"},
	)
	prometheus.MustRegister(total)
	prometheus.MustRegister(values)

	myTypesMetrics = PrometheusMetricsTypes{
		TotalTasks:  total,
		TotalValues: values,
	}
}

func IncreaseProcessedTasks() {
	myGeneralMetrics.ProcessedTasks.Inc()
}

func IncreaseDoneTasks() {
	myGeneralMetrics.ProcessedTasks.Dec()
	myGeneralMetrics.DoneTasks.Inc()
}

func IncreaseTotalTasksAndValue(taskType int8, taskValue int8) {
	myTypesMetrics.TotalTasks.With(prometheus.Labels{"type": fmt.Sprintf("type %d", taskType)}).Inc()
	myTypesMetrics.TotalValues.With(prometheus.Labels{"type": fmt.Sprintf("type %d", taskType)}).Add(float64(taskValue))
}

func StartPrometheusServer(addr string) {
	http.Handle("/metrics", promhttp.Handler()) // /metrics endpoint for Prometheus
	log.Println("Prometheus metrics available at /metrics)")
	log.Fatal(http.ListenAndServe(addr, nil)) // Prometheus listens on port "addr"
}
