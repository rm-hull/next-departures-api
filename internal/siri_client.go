package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/rm-hull/next-departures-api/internal/models/siri"
)

type SiriClient struct {
	appId    string
	endpoint string
}

func NewSiriClient(appId, appKey string) *SiriClient {
	return &SiriClient{
		appId:    appId,
		endpoint: fmt.Sprintf("https://%s:%s@transportapi.com/nextbuses", appId, appKey),
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
				PreviewInterval:   "PT120M",
				MaximumStopVisits: 10,
			},
		},
	}

	body, err := xml.Marshal(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal XML request: %w", err)
	}
	resp, err := http.Post(c.endpoint, "application/xml", bytes.NewReader(body))
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
