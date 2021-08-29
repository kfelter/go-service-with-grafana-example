// this example is based off this blog post
// https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

func main() {
	// register our custom metrics
	prometheus.Register(httpDuration)
	prometheus.Register(totalRequests)

	// Serving static files
	http.HandleFunc("/example", func(rw http.ResponseWriter, r *http.Request) {
		// write response header
		rw.WriteHeader(200)

		// set up observation to record response times of this endpoint
		timer := prometheus.NewTimer(httpDuration.WithLabelValues("/example"))
		defer timer.ObserveDuration()

		sleepms := time.Duration(rand.Intn(500)) * time.Millisecond
		time.Sleep(sleepms)

		rw.Write([]byte(fmt.Sprintf(`hello from go %v`, sleepms)))

		// increment page views
		totalRequests.WithLabelValues("/example").Inc()

	})

	// Prometheus endpoint
	// this will display some information about the go runtime and our custom metrics
	http.Handle("/prometheus", promhttp.Handler())

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", nil)
	log.Fatal(err)
}
