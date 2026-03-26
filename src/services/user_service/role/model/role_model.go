package model

import "url_shortener_pro/src/common/helper"

type Role struct {
	Id          int64             `gorm:"primaryKey"`
	RoleName    string            `gorm:"unique;not null;column:role_name"`
	Permissions helper.JsonObject `gorm:"jsonb"`
}

type Permission struct {
	Id           int64  `gorm:"primaryKey"`
	Path         string `gorm:"not null;column:path"`
	Method       string `gorm:"not null;column:method"`
	EndpointName string `gorm:"not null;column:endpoint_name"`
}

type UserRole struct {
	Id     int64 `json:"id" gorm:"primaryKey"`
	UserId int64 `json:"userId" gorm:"not null;column:user_id"`
	RoleId int64 `json:"roleId" gorm:"not null;column:role_id"`
}
