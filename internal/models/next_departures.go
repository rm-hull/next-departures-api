package models

import "time"

type NextDepartureResponse struct {
	Results     []NextDeparture `json:"results"`
	Attribution []string        `json:"attribution"`
}

type NextDeparture struct {
	LineName              string     `json:"line_name"`
	Destination           string     `json:"destination"`
	OperatorRef           string     `json:"operator_ref"`
	AimedDepartureTime    *time.Time `json:"aimed_departure_time,omitempty"`
	ExpectedDepartureTime *time.Time `json:"expected_departure_time,omitempty"`
}
