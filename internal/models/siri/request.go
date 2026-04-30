package siri

import (
	"encoding/xml"
)

type StopMonitoringRequest struct {
	XMLName        xml.Name       `xml:"Siri"`
	Version        string         `xml:"version,attr"`
	Xmlns          string         `xml:"xmlns,attr"`
	ServiceRequest ServiceRequest `xml:"ServiceRequest"`
}

type ServiceRequest struct {
	RequestorRef          string            `xml:"RequestorRef"`
	StopMonitoringRequest StopMonitoringReq `xml:"StopMonitoringRequest"`
}

type StopMonitoringReq struct {
	Version           string `xml:"version,attr"`
	MonitoringRef     string `xml:"MonitoringRef"`
	PreviewInterval   string `xml:"PreviewInterval"` // e.g. "PT30M"
	MaximumStopVisits int    `xml:"MaximumStopVisits,omitempty"`
}
