package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rm-hull/next-departures-api/internal"
	"github.com/rm-hull/next-departures-api/internal/models"
)

func NextDepartures(client *internal.SiriClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		stopId := c.Param("stopId")
		siri, statusCode, err := client.GetStopMonitoring(stopId)
		if err != nil {
			log.Printf("error while fetching next departures: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			return
		}

		switch statusCode {
		case http.StatusOK:
			departures := make([]models.NextDeparture, 0)
			if len(siri.ServiceDelivery.StopMonitoringDelivery) == 0 {
				c.JSON(http.StatusOK, models.NextDepartureResponse{
					Results:     departures,
					Attribution: internal.ATTRIBUTION,
				})
				return
			}
			for _, visit := range siri.ServiceDelivery.StopMonitoringDelivery[0].MonitoredStopVisit {
				departures = append(departures, models.NextDeparture{
					LineName:              visit.MonitoredVehicleJourney.PublishedLineName,
					Destination:           visit.MonitoredVehicleJourney.DirectionName,
					OperatorRef:           visit.MonitoredVehicleJourney.OperatorRef,
					AimedDepartureTime:    visit.MonitoredVehicleJourney.MonitoredCall.AimedDepartureTime,
					ExpectedDepartureTime: visit.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime,
				})
			}

			c.JSON(http.StatusOK, models.NextDepartureResponse{
				Results:     departures,
				Attribution: internal.ATTRIBUTION,
			})

		case http.StatusBadRequest:
			errMsg := "Bad request to SIRI API"
			if siri.ServiceDelivery.ErrorCondition != nil && siri.ServiceDelivery.ErrorCondition.OtherError != nil {
				errMsg = siri.ServiceDelivery.ErrorCondition.OtherError.ErrorText
			}
			c.JSON(statusCode, gin.H{"error": errMsg})

		case http.StatusForbidden, http.StatusUnauthorized:
			errMsg := "Access denied"
			if siri.ServiceDelivery.ErrorCondition != nil && siri.ServiceDelivery.ErrorCondition.AccessNotAllowedError != nil {
				errMsg = siri.ServiceDelivery.ErrorCondition.AccessNotAllowedError.ErrorText
			}
			log.Printf("unexpected HTTP status code (%d) from SIRI API: %s", statusCode, errMsg)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
			
		default:
			log.Printf("unexpected HTTP status code (%d) from SIRI API", statusCode)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal server error occurred"})
		}
	}
}
