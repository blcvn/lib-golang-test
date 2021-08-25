package main

import (
	"net/http"

	"log"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Counters are a simple metric type that can only be incremented or be reset to zero on restart.
	// It is often used to count primitive data like the total number of requests to a services or number of tasks completed
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "my_counter",
			Help:      "This is my counter",
		})

	//Gauges also represent a single numerical value but different to counters the value can go up as well as down.
	//Therefore gauges are often used for measured values like temperature, humidy or current memory usage.
	gauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	//Histograms are used to measure the frequency of value observations that fall into specific predefined buckets
	//This means that they will provide information about the distribution of a metric like response time and signal outliers
	//By default Prometheus provides the following buckets: .005, .01, .025, .05, .075, .1, .25, .5, .75, 1, 2.5, 5, 7.5, 10

	histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	// Summaries are very similar to Histograms because they both expose the distribution of a given data set.
	// The one major difference is that a Histogram estimate quantiles on the Prometheus server while Summaries are calculated on the client side

	summary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)

func main() {
	rand.Seed(time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(counter)
	prometheus.MustRegister(gauge)
	prometheus.MustRegister(histogram)
	prometheus.MustRegister(summary)

	go func() {
		for {
			counter.Add(rand.Float64() * 5)
			gauge.Add(rand.Float64()*15 - 5)
			histogram.Observe(rand.Float64() * 10)
			summary.Observe(rand.Float64() * 10)

			time.Sleep(time.Second)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
