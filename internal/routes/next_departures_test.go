package routes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rm-hull/next-departures-api/internal/models/siri"
	"github.com/stretchr/testify/assert"
)

type mockSiriClient struct {
	getStopMonitoringFn func(monitoringRef string) (*siri.Siri, int, error)
}

func (m *mockSiriClient) GetStopMonitoring(monitoringRef string) (*siri.Siri, int, error) {
	return m.getStopMonitoringFn(monitoringRef)
}

func TestNextDepartures_RateLimited(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return &siri.Siri{
				ServiceDelivery: siri.ServiceDelivery{
					ErrorCondition: &siri.ErrorCondition{
						AccessNotAllowedError: &siri.Error{
							ErrorText: "Usage limits are exceeded",
						},
					},
				},
			}, http.StatusForbidden, nil
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "Rate limit exceeded. Please try again after midnight.")
	assert.NotEmpty(t, w.Header().Get("Retry-After"))
	assert.NotEmpty(t, w.Header().Get("X-RateLimit-Reset"))
}

func TestNextDepartures_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	aimedTime := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	expectedTime := time.Date(2026, 5, 1, 12, 5, 0, 0, time.UTC)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return &siri.Siri{
				ServiceDelivery: siri.ServiceDelivery{
					StopMonitoringDelivery: []siri.StopMonitoringDelivery{
						{
							MonitoredStopVisit: []siri.MonitoredStopVisit{
								{
									MonitoredVehicleJourney: siri.MonitoredVehicleJourney{
										PublishedLineName: "42",
										DirectionName:     "Galaxy",
										OperatorRef:       "MARVIN",
										MonitoredCall: siri.MonitoredCall{
											AimedDepartureTime:    &aimedTime,
											ExpectedDepartureTime: &expectedTime,
										},
									},
								},
							},
						},
					},
				},
			}, http.StatusOK, nil
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"line_name":"42"`)
	assert.Contains(t, w.Body.String(), `"destination":"Galaxy"`)
	assert.Contains(t, w.Body.String(), `"attribution"`)
}

func TestNextDepartures_NoDepartures(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return &siri.Siri{
				ServiceDelivery: siri.ServiceDelivery{
					StopMonitoringDelivery: []siri.StopMonitoringDelivery{},
				},
			}, http.StatusOK, nil
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"results":[]`)
}

func TestNextDepartures_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return &siri.Siri{
				ServiceDelivery: siri.ServiceDelivery{
					ErrorCondition: &siri.ErrorCondition{
						OtherError: &siri.Error{
							ErrorText: "Invalid monitoring reference",
						},
					},
				},
			}, http.StatusBadRequest, nil
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid monitoring reference")
}

func TestNextDepartures_AccessDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return &siri.Siri{
				ServiceDelivery: siri.ServiceDelivery{
					ErrorCondition: &siri.ErrorCondition{
						AccessNotAllowedError: &siri.Error{
							ErrorText: "Invalid API Key",
						},
					},
				},
			}, http.StatusUnauthorized, nil
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// In the current implementation, StatusUnauthorized/StatusForbidden (not rate-limited)
	// results in an Internal Server Error being returned to the caller.
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "An internal server error occurred")
}

func TestNextDepartures_ClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockClient := &mockSiriClient{
		getStopMonitoringFn: func(monitoringRef string) (*siri.Siri, int, error) {
			return nil, 0, errors.New("network failure")
		},
	}

	r := gin.New()
	r.GET("/v1/next-departures/:stopId", NextDepartures(mockClient))

	req, _ := http.NewRequest(http.MethodGet, "/v1/next-departures/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "An internal server error occurred")
}
