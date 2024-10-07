package metrics

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/nicoletavoinea/GolangProducerConsumer/internal/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var myGeneralMetrics PrometheusMetricsGeneral
var myTypesMetrics PrometheusMetricsTypes
var mu sync.Mutex
var mu1 sync.Mutex
var mu2 sync.Mutex

type PrometheusMetricsGeneral struct {
	ProcessedTasks prometheus.Gauge
	DoneTasks      prometheus.Gauge
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

	done := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "done_tasks",
			Help: "Total number of done tasks.",
		},
	)

	//set to initial values, computed from database
	processed.Set(database.GetNumberOfProcessingTasks())
	done.Set(database.GetNumberOfDoneTasks())

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

	totalNumbers := database.GetNumberOfTasksByType()
	for i := 0; i < 10; i++ {
		total.With(prometheus.Labels{"type": fmt.Sprintf("type %d", i)}).Add(float64(totalNumbers[i]))
	}

	totalValues := database.GetValueOfTasksByType()
	for i := 0; i < 10; i++ {
		values.With(prometheus.Labels{"type": fmt.Sprintf("type %d", i)}).Add(float64(totalValues[i]))
	}

	prometheus.MustRegister(total)
	prometheus.MustRegister(values)

	myTypesMetrics = PrometheusMetricsTypes{
		TotalTasks:  total,
		TotalValues: values,
	}
}

func IncreaseProcessedTasks() { //increment number of processing tasks
	mu.Lock()
	myGeneralMetrics.ProcessedTasks.Inc()
	mu.Unlock()
}

func IncreaseDoneTasks() { //increment number of done tasks
	mu.Lock()
	myGeneralMetrics.ProcessedTasks.Dec()
	mu.Unlock()
	mu1.Lock()
	myGeneralMetrics.DoneTasks.Inc()
	mu1.Unlock()
}

func IncreaseTotalTasksAndValue(taskType int32, taskValue int32) { //increment number of tasks per task type & their values
	mu2.Lock()
	myTypesMetrics.TotalTasks.With(prometheus.Labels{"type": fmt.Sprintf("type %d", taskType)}).Inc()
	myTypesMetrics.TotalValues.With(prometheus.Labels{"type": fmt.Sprintf("type %d", taskType)}).Add(float64(taskValue))
	mu2.Unlock()
}

func StartPrometheusServer(addr string) { //startprometheus server at specified addr
	http.Handle("/metrics", promhttp.Handler()) // /metrics endpoint for Prometheus
	log.Printf("Prometheus metrics available at %v/metrics)", addr)
	log.Fatal(http.ListenAndServe(addr, nil)) // Prometheus listens on port "addr"
}
