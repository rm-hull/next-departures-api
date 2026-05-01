package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestSiriClient_GetStopMonitoring_Metrics(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Siri xmlns="http://www.siri.org.uk/siri">
    <ServiceDelivery>
        <ResponseTimestamp>2026-05-01T12:00:00Z</ResponseTimestamp>
        <Status>true</Status>
        <StopMonitoringDelivery version="1.0">
            <ResponseTimestamp>2026-05-01T12:00:00Z</ResponseTimestamp>
            <MonitoredStopVisit>
                <RecordedAtTime>2026-05-01T11:59:00Z</RecordedAtTime>
            </MonitoredStopVisit>
            <MonitoredStopVisit>
                <RecordedAtTime>2026-05-01T11:59:30Z</RecordedAtTime>
            </MonitoredStopVisit>
        </StopMonitoringDelivery>
    </ServiceDelivery>
</Siri>`)
	}))
	defer ts.Close()

	reg := prometheus.NewRegistry()
	client := NewSiriClient("appId", "appKey", reg).(*siriClient)
	client.endpoint = ts.URL // Override endpoint for test

	_, statusCode, err := client.GetStopMonitoring("123")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	// Verify metrics
	metricFamilies, err := reg.Gather()
	assert.NoError(t, err)

	foundLatency := false
	foundStatusCodes := false
	foundItemsFetched := false

	for _, mf := range metricFamilies {
		switch mf.GetName() {
		case "next_departures_siri_api_http_response_latency_seconds":
			foundLatency = true
		case "next_departures_siri_api_http_response_status_codes_total":
			foundStatusCodes = true
			assert.Equal(t, float64(1), mf.GetMetric()[0].GetCounter().GetValue())
		case "next_departures_siri_items_fetched_total":
			foundItemsFetched = true
			assert.Equal(t, float64(2), mf.GetMetric()[0].GetCounter().GetValue())
		}
	}

	assert.True(t, foundLatency)
	assert.True(t, foundStatusCodes)
	assert.True(t, foundItemsFetched)
}
