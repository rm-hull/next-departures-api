package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/rm-hull/next-departures-api/internal/models/siri"
)

type SiriClient struct {
	appId             string
	appKey            string
	endpoint          string
	previewInterval   string
	maximumStopVisits int
	httpClient        *http.Client
}

func NewSiriClient(appId, appKey string) *SiriClient {
	return &SiriClient{
		appId:             appId,
		appKey:            appKey,
		endpoint:          "https://transportapi.com/nextbuses",
		previewInterval:   "PT120M",
		maximumStopVisits: 10,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *SiriClient) GetStopMonitoring(monitoringRef string) (*siri.Siri, int, error) {
	req := siri.StopMonitoringRequest{
		Version: "1.0",
		Xmlns:   "http://www.siri.org.uk/siri",
		ServiceRequest: siri.ServiceRequest{
			RequestorRef: c.appId,
			StopMonitoringRequest: siri.StopMonitoringReq{
				Version:           "1.0",
				MonitoringRef:     monitoringRef,
				PreviewInterval:   c.previewInterval,
				MaximumStopVisits: c.maximumStopVisits,
			},
		},
	}

	body, err := xml.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal XML request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/xml")
	httpReq.SetBasicAuth(c.appId, c.appKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var siri siri.Siri
	err = xml.NewDecoder(resp.Body).Decode(&siri)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to decode XML response: %w", err)
	}

	return &siri, resp.StatusCode, nil
}
