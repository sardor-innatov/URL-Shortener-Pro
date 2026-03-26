package dto

type StatsDto struct {
	TotalClicks   int64             `json:"totalClicks"`
	UniqueIps     int64             `json:"uniqueIps"`
	TopCountries  []CountryStatsDto `json:"topCountries"`
	Last24hClicks int               `json:"last24hClicks"`
}

type CountryStatsDto struct {
	Country string `json:"country" gorm:"country"`
	Count   int64  `json:"count" gorm:"count"`
}
