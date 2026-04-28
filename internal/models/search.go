package models

import "time"

type SearchResult struct {
	NaPTAN
}

type SearchResponse struct {
	Results     []SearchResult `json:"results"`
	Attribution []string       `json:"attribution"`
	LastUpdated *time.Time     `json:"last_updated,omitempty"`
}
