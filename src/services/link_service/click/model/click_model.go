package model

import "time"

type Click struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	LinkId    int64     `json:"linkId" gorm:"column:link_id;index:idx_link_country;index:idx_clicks_link_id_created_at"`
	IpAddress string    `json:"ipAddress" gorm:"column:ip_address"`
	UserAgent string    `json:"userAgent" gorm:"column:user_agent"`
	Country   *string   `json:"country" gorm:"column:country;index:idx_link_country"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime;index:idx_clicks_link_id_created_at"`
}
