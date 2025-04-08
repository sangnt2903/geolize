package model

type IPUpdateRequest struct {
	IP                 string              `json:"ip"`
	Continent          *Continent          `json:"continent"`
	Country            *Country            `json:"country"`
	Location           *Location           `json:"location"`
	Subdivisions       []*Subdivision      `json:"subdivisions"`
	Postal             *Postal             `json:"postal"`
	City               *City               `json:"city"`
	RepresentedCountry *RepresentedCountry `json:"represented_country"`
	RegisteredCountry  *RegisteredCountry  `json:"registered_country"`
	Traits             *Traits             `json:"traits"`
}
