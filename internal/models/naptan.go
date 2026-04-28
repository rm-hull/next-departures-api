package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const NAPTAN_CSV_URL = "https://naptan.api.dft.gov.uk/v1/access-nodes?dataFormat=csv"

type NaPTAN struct {
	ATCOCode                string     `json:"atco_code,omitempty"`
	NaptanCode              string     `json:"naptan_code,omitempty"`
	PlateCode               string     `json:"plate_code,omitempty"`
	CleardownCode           string     `json:"cleardown_code,omitempty"`
	CommonName              string     `json:"common_name,omitempty"`
	CommonNameLang          string     `json:"-"`
	ShortCommonName         string     `json:"short_common_name,omitempty"`
	ShortCommonNameLang     string     `json:"-"`
	Landmark                string     `json:"landmark,omitempty"`
	LandmarkLang            string     `json:"-"`
	Street                  string     `json:"street,omitempty"`
	StreetLang              string     `json:"-"`
	Crossing                string     `json:"crossing,omitempty"`
	CrossingLang            string     `json:"-"`
	Indicator               string     `json:"indicator,omitempty"`
	IndicatorLang           string     `json:"-"`
	Bearing                 string     `json:"bearing,omitempty"`
	NptgLocalityCode        string     `json:"nptg_locality_code,omitempty"`
	LocalityName            string     `json:"locality_name,omitempty"`
	ParentLocalityName      string     `json:"parent_locality_name,omitempty"`
	GrandParentLocalityName string     `json:"grand_parent_locality_name,omitempty"`
	Town                    string     `json:"town,omitempty"`
	TownLang                string     `json:"-"`
	Suburb                  string     `json:"suburb,omitempty"`
	SuburbLang              string     `json:"-"`
	LocalityCentre          *bool      `json:"locality_centre,omitempty"`
	GridType                string     `json:"grid_type,omitempty"`
	Easting                 *int       `json:"easting,omitempty"`
	Northing                *int       `json:"northing,omitempty"`
	Longitude               *float64   `json:"longitude,omitempty"`
	Latitude                *float64   `json:"latitude,omitempty"`
	StopType                string     `json:"stop_type,omitempty"`
	BusStopType             string     `json:"bus_stop_type,omitempty"`
	TimingStatus            string     `json:"timing_status,omitempty"`
	DefaultWaitTime         string     `json:"default_wait_time,omitempty"`
	Notes                   string     `json:"notes,omitempty"`
	NotesLang               string     `json:"-"`
	AdministrativeAreaCode  string     `json:"administrative_area_code,omitempty"`
	CreationDateTime        *time.Time `json:"creation_date_time,omitempty"`
	ModificationDateTime    *time.Time `json:"modification_date_time,omitempty"`
	RevisionNumber          *int       `json:"revision_number,omitempty"`
	Modification            string     `json:"modification,omitempty"`
	Status                  string     `json:"status,omitempty"`
}

func FromTuple(data []string, headers []string) (*NaPTAN, error) {
	localityCentre, err := parseBool(data[25])
	if err != nil {
		return nil, err
	}
	easting, err := parseInt(data[27])
	if err != nil {
		return nil, err
	}
	northing, err := parseInt(data[28])
	if err != nil {
		return nil, err
	}
	longitude, err := parseFloat(data[29])
	if err != nil {
		return nil, err
	}
	latitude, err := parseFloat(data[30])
	if err != nil {
		return nil, err
	}
	creationDateTime, err := parseTime(data[38])
	if err != nil {
		return nil, err
	}
	modificationDateTime, err := parseTime(data[39])
	if err != nil {
		return nil, err
	}
	revisionNumber, err := parseInt(data[40])
	if err != nil {
		return nil, err
	}

	return &NaPTAN{
		ATCOCode:                data[0],
		NaptanCode:              data[1],
		PlateCode:               data[2],
		CleardownCode:           data[3],
		CommonName:              data[4],
		CommonNameLang:          data[5],
		ShortCommonName:         data[6],
		ShortCommonNameLang:     data[7],
		Landmark:                data[8],
		LandmarkLang:            data[9],
		Street:                  data[10],
		StreetLang:              data[11],
		Crossing:                data[12],
		CrossingLang:            data[13],
		Indicator:               data[14],
		IndicatorLang:           data[15],
		Bearing:                 data[16],
		NptgLocalityCode:        data[17],
		LocalityName:            data[18],
		ParentLocalityName:      data[19],
		GrandParentLocalityName: data[20],
		Town:                    data[21],
		TownLang:                data[22],
		Suburb:                  data[23],
		SuburbLang:              data[24],
		LocalityCentre:          localityCentre,
		GridType:                data[26],
		Easting:                 easting,
		Northing:                northing,
		Longitude:               longitude,
		Latitude:                latitude,
		StopType:                data[31],
		BusStopType:             data[32],
		TimingStatus:            data[33],
		DefaultWaitTime:         data[34],
		Notes:                   data[35],
		NotesLang:               data[36],
		AdministrativeAreaCode:  data[37],
		CreationDateTime:        creationDateTime,
		ModificationDateTime:    modificationDateTime,
		RevisionNumber:          revisionNumber,
		Modification:            data[41],
		Status:                  data[42],
	}, nil
}

func (n *NaPTAN) ToTuple() []any {
	strOrNil := func(s string) any {
		if s == "" {
			return nil
		}
		return s
	}

	derefInt := func(i *int) any {
		if i == nil {
			return nil
		}
		return *i
	}

	derefBool := func(b *bool) any {
		if b == nil {
			return nil
		}
		return *b
	}

	derefFloat := func(f *float64) any {
		if f == nil {
			return nil
		}
		return *f
	}

	return []any{
		strOrNil(n.ATCOCode),
		strOrNil(n.NaptanCode),
		strOrNil(n.PlateCode),
		strOrNil(n.CleardownCode),
		strOrNil(n.CommonName),
		strOrNil(n.CommonNameLang),
		strOrNil(n.ShortCommonName),
		strOrNil(n.ShortCommonNameLang),
		strOrNil(n.Landmark),
		strOrNil(n.LandmarkLang),
		strOrNil(n.Street),
		strOrNil(n.StreetLang),
		strOrNil(n.Crossing),
		strOrNil(n.CrossingLang),
		strOrNil(n.Indicator),
		strOrNil(n.IndicatorLang),
		strOrNil(n.Bearing),
		strOrNil(n.NptgLocalityCode),
		strOrNil(n.LocalityName),
		strOrNil(n.ParentLocalityName),
		strOrNil(n.GrandParentLocalityName),
		strOrNil(n.Town),
		strOrNil(n.TownLang),
		strOrNil(n.Suburb),
		strOrNil(n.SuburbLang),
		derefBool(n.LocalityCentre),
		strOrNil(n.GridType),
		derefInt(n.Easting),
		derefInt(n.Northing),
		derefFloat(n.Longitude),
		derefFloat(n.Latitude),
		strOrNil(n.StopType),
		strOrNil(n.BusStopType),
		strOrNil(n.TimingStatus),
		strOrNil(n.DefaultWaitTime),
		strOrNil(n.Notes),
		strOrNil(n.NotesLang),
		strOrNil(n.AdministrativeAreaCode),
		n.CreationDateTime,
		n.ModificationDateTime,
		derefInt(n.RevisionNumber),
		strOrNil(n.Modification),
		strOrNil(n.Status),
	}
}

func parseInt(s string) (*int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int from '%s': %w", s, err)
	}
	return &i, nil
}

func parseBool(s string) (*bool, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bool from '%s': %w", s, err)
	}
	return &b, nil
}

func parseFloat(s string) (*float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse float from '%s': %w", s, err)
	}
	return &f, nil
}

func parseTime(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	// Try parsing in order of most specific to least specific formats
	var t time.Time
	var err error

	// First try RFC3339Nano (with milliseconds and timezone)
	t, err = time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return &t, nil
	}

	// Then try RFC3339 (with timezone, no milliseconds)
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return &t, nil
	}

	// Finally try simple format (no timezone, no milliseconds)
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return &t, nil
	}

	return nil, fmt.Errorf("failed to parse time from '%s' with any supported format: %w", s, err)
}
