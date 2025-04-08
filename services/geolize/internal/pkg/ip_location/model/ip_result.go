package model

type IPResult struct {
	IP                 string              `json:"ip,omitempty"`
	DBVersion          string              `json:"db_version,omitempty"`
	Continent          *Continent          `json:"continent,omitempty"`
	Country            *Country            `json:"country,omitempty"`
	Location           *Location           `json:"location,omitempty"`
	Subdivisions       []*Subdivision      `json:"subdivisions,omitempty"`
	Postal             *Postal             `json:"postal,omitempty"`
	City               *City               `json:"city,omitempty"`
	RepresentedCountry *RepresentedCountry `json:"represented_country,omitempty"`
	RegisteredCountry  *RegisteredCountry  `json:"registered_country,omitempty"`
	Traits             *Traits             `json:"traits,omitempty"`
}

type Continent struct {
	Code  string            `json:"code,omitempty"`
	Names map[string]string `json:"names,omitempty"`
}

type Country struct {
	ISOCode           string            `json:"iso_code,omitempty"`
	Names             map[string]string `json:"names,omitempty"`
	IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
}

type Location struct {
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	AccuracyRadius uint16  `json:"accuracy_radius,omitempty"`
	TimeZone       string  `json:"time_zone,omitempty"`
}

type Subdivision struct {
	ISOCode string            `json:"iso_code,omitempty"`
	Names   map[string]string `json:"names,omitempty"`
}

type Postal struct {
	Code string `json:"code,omitempty"`
}

type City struct {
	Names map[string]string `json:"names,omitempty"`
}

type RepresentedCountry struct {
	ISOCode           string            `json:"iso_code,omitempty"`
	Names             map[string]string `json:"names,omitempty"`
	Type              string            `json:"type,omitempty"`
	IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
}

type RegisteredCountry struct {
	ISOCode           string            `json:"iso_code,omitempty"`
	Names             map[string]string `json:"names,omitempty"`
	IsInEuropeanUnion bool              `json:"is_in_european_union,omitempty"`
}

type Traits struct {
	IsAnonymousProxy    bool `json:"is_anonymous_proxy,omitempty"`
	IsAnycast           bool `json:"is_anycast,omitempty"`
	IsSatelliteProvider bool `json:"is_satellite_provider,omitempty"`
}
