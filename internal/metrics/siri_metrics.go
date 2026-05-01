package metrics

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	responseLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "next_departures_siri_api_http_response_latency_seconds",
			Help:    "TransportAPI SIRI client HTTP response latency in seconds.",
			Buckets: []float64{0.1, 0.25, 0.5, 1, 2, 5, 10, 30, math.Inf(1)},
		},
		[]string{"method"},
	)
	responseStatusCode = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "next_departures_siri_api_http_response_status_codes_total",
			Help: "TransportAPI SIRI client total number of HTTP responses by status code.",
		},
		[]string{"method", "status_code"},
	)
	itemsFetchedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "next_departures_siri_items_fetched_total",
			Help: "TransportAPI SIRI client total number of stop visits fetched.",
		},
		[]string{"method"},
	)
)

type SiriMetrics struct {
	ResponseLatency    *prometheus.HistogramVec
	ResponseStatusCode *prometheus.CounterVec
	ItemsFetchedTotal  *prometheus.CounterVec
}

func NewSiriMetrics(reg prometheus.Registerer) *SiriMetrics {
	RegisterOrPanic(reg,
		responseLatency,
		responseStatusCode,
		itemsFetchedTotal,
	)

	return &SiriMetrics{
		ResponseLatency:    responseLatency,
		ResponseStatusCode: responseStatusCode,
		ItemsFetchedTotal:  itemsFetchedTotal,
	}
}

func (m *SiriMetrics) RecordHttpCall(start time.Time, method string, resp *http.Response, err error) {
	if m == nil {
		return
	}
	m.ResponseLatency.WithLabelValues(method).Observe(time.Since(start).Seconds())
	if err == nil && resp != nil {
		m.ResponseStatusCode.WithLabelValues(method, strconv.Itoa(resp.StatusCode)).Inc()
	} else if err != nil {
		m.ResponseStatusCode.WithLabelValues(method, "error").Inc()
	}
}

func (m *SiriMetrics) RecordFetchedItems(method string, count int) {
	if m == nil {
		return
	}
	if count > 0 {
		m.ItemsFetchedTotal.WithLabelValues(method).Add(float64(count))
	}
}
