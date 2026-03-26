package model

import (
	"time"
)

type Link struct {
	Id        int64      `json:"id" gorm:"primaryKey"`
	UserId    int64      `json:"userId" gorm:"user_id"`
	URL       string     `json:"url" gorm:"column:url"`
	ShortCode string     `json:"shortCode" gorm:"type:varchar(32);column:short_code;index;unique;not null;"`
	IsDeleted bool       `json:"isDeleted" gorm:"column:is_deleted;default:false"`
	ExpiresAt *time.Time `json:"expiresAt" gorm:"column:expires_at;default:NULL"`
	CreatedAt *time.Time `json:"createdAt" gorm:"column:created_at"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}
