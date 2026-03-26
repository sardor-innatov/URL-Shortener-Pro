package dto

import "time"

type LinkCreateDto struct {
	URL         string `json:"url" validate:"required,http_url"`
	CustomAlias string `json:"customAlias" validate:"omitempty,min=4,max=32,hostname_rfc1123"`
	ExpiresIn   int    `json:"expiresIn"`
}

type RedirectDto struct {
	ShortCode string `json:"shortCode" gorm:"column:short_code;index;unique;not null"`
}

type LinkRedisDto struct {
	URL       string     `json:"url"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

type LinkGetDto struct {
	URL       string `json:"url" gorm:"url"`
	ShortCode string `json:"shortCode" gorm:"short_code"`
	ExpiresAt *time.Time    `json:"expiresAt" gorm:"expires_at"`
}
