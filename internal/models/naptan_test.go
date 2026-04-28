package models

import (
	"testing"
)

func TestFromTuple_DynamicHeaders(t *testing.T) {
	headers := []string{"ATCOCode", "NaptanCode", "CommonName", "LocalityCentre", "Easting", "Northing", "Longitude", "Latitude"}
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	data := make([]string, len(headers))
	data[headerMap["ATCOCode"]] = "0100BRP90311"
	data[headerMap["NaptanCode"]] = "bstgwpm"
	data[headerMap["CommonName"]] = "Temple Meads Stn"
	data[headerMap["LocalityCentre"]] = "true"
	data[headerMap["Easting"]] = "359403"
	data[headerMap["Northing"]] = "172513"
	data[headerMap["Longitude"]] = "-2.5856"
	data[headerMap["Latitude"]] = "51.45014"

	naptan, err := FromTuple(data, headerMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if naptan.ATCOCode != "0100BRP90311" {
		t.Errorf("expected ATCOCode 0100BRP90311, got %s", naptan.ATCOCode)
	}
	if naptan.NaptanCode != "bstgwpm" {
		t.Errorf("expected NaptanCode bstgwpm, got %s", naptan.NaptanCode)
	}
	if naptan.CommonName != "Temple Meads Stn" {
		t.Errorf("expected CommonName Temple Meads Stn, got %s", naptan.CommonName)
	}
	if naptan.LocalityCentre == nil || !*naptan.LocalityCentre {
		t.Errorf("expected LocalityCentre true, got %v", naptan.LocalityCentre)
	}
	if naptan.Easting == nil || *naptan.Easting != 359403 {
		t.Errorf("expected Easting 359403, got %v", naptan.Easting)
	}
	if naptan.Northing == nil || *naptan.Northing != 172513 {
		t.Errorf("expected Northing 172513, got %v", naptan.Northing)
	}
	if naptan.Longitude == nil || *naptan.Longitude != -2.5856 {
		t.Errorf("expected Longitude -2.5856, got %v", naptan.Longitude)
	}
	if naptan.Latitude == nil || *naptan.Latitude != 51.45014 {
		t.Errorf("expected Latitude 51.45014, got %v", naptan.Latitude)
	}
}

func TestFromTuple_ReorderedHeaders(t *testing.T) {
	// Different order than standard
	headers := []string{"CommonName", "ATCOCode", "LocalityCentre"}
	headerMap := make(map[string]int)
	for i, h := range headers {
		headerMap[h] = i
	}

	data := []string{"My Stop", "12345", "false"}

	naptan, err := FromTuple(data, headerMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if naptan.CommonName != "My Stop" {
		t.Errorf("expected CommonName My Stop, got %s", naptan.CommonName)
	}
	if naptan.ATCOCode != "12345" {
		t.Errorf("expected ATCOCode 12345, got %s", naptan.ATCOCode)
	}
	if naptan.LocalityCentre == nil || *naptan.LocalityCentre {
		t.Errorf("expected LocalityCentre false, got %v", naptan.LocalityCentre)
	}
}
