package model

import "time"

type User struct {
	Id        int64  `json:"id" gorm:"primaryKey"`
	FirstName string `json:"firstName" gorm:"column:first_name"`
	LastName  string `json:"lastName" gorm:"column:last_name"`
	Email     string `json:"email" gorm:"column:email;unique"`
	Password  string `json:"password" gorm:"column:password"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}