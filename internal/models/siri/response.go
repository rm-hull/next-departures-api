package siri

import (
	"encoding/xml"
	"time"
)

// Top-level envelope
type Siri struct {
	XMLName         xml.Name        `xml:"Siri"`
	Xmlns           string          `xml:"xmlns,attr"`
	ServiceDelivery ServiceDelivery `xml:"ServiceDelivery"`
}

type ServiceDelivery struct {
	ResponseTimestamp      time.Time                `xml:"ResponseTimestamp"`
	StopMonitoringDelivery []StopMonitoringDelivery `xml:"StopMonitoringDelivery"`
	Status                 bool                     `xml:"Status"`
	ErrorCondition         *ErrorCondition          `xml:"ErrorCondition,omitempty"`
}

type ErrorCondition struct {
	AccessNotAllowedError *Error `xml:"AccessNotAllowedError,omitempty"`
	OtherError            *Error `xml:"OtherError,omitempty"`
}

type Error struct {
	ErrorText string `xml:"ErrorText"`
}

type StopMonitoringDelivery struct {
	Version            string               `xml:"version,attr"`
	ResponseTimestamp  time.Time            `xml:"ResponseTimestamp"`
	ValidUntil         time.Time            `xml:"ValidUntil"`
	MonitoredStopVisit []MonitoredStopVisit `xml:"MonitoredStopVisit"`
}

type MonitoredStopVisit struct {
	RecordedAtTime          time.Time               `xml:"RecordedAtTime"`
	MonitoringRef           string                  `xml:"MonitoringRef"`
	MonitoredVehicleJourney MonitoredVehicleJourney `xml:"MonitoredVehicleJourney"`
	Extensions              Extensions              `xml:"Extensions"`
}

type MonitoredVehicleJourney struct {
	VehicleMode       string        `xml:"VehicleMode"`
	PublishedLineName string        `xml:"PublishedLineName"`
	DirectionName     string        `xml:"DirectionName"`
	OperatorRef       string        `xml:"OperatorRef"`
	MonitoredCall     MonitoredCall `xml:"MonitoredCall"`
}

type MonitoredCall struct {
	AimedArrivalTime      *time.Time `xml:"AimedArrivalTime"`
	AimedDepartureTime    *time.Time `xml:"AimedDepartureTime"`
	ExpectedArrivalTime   *time.Time `xml:"ExpectedArrivalTime"`
	ExpectedDepartureTime *time.Time `xml:"ExpectedDepartureTime"`
}

type Extensions struct {
	VehicleJourney VehicleJourney `xml:"VehicleJourney"`
}

type VehicleJourney struct {
	SeatedOccupancy     int `xml:"SeatedOccupancy"`
	SeatedCapacity      int `xml:"SeatedCapacity"`
	WheelchairOccupancy int `xml:"WheelchairOccupancy"`
	WheelchairCapacity  int `xml:"WheelchairCapacity"`
}
