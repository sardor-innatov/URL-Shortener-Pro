package model

import "time"

type LinkStats struct {
	LinkID        uint `gorm:"primaryKey"`
	TotalClicks   int64  `gorm:"default:0"`
	UniqueIPs     int  `gorm:"default:0"`
	Last24hClicks int  `gorm:"default:0;column:last_24h_clicks"`
	UpdatedAt     time.Time
}
