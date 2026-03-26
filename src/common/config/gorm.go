package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type gormInstance struct {
	*gorm.DB
}

var _gormInstance *gormInstance

func NewGorm(c *Config) (*gorm.DB, error) {

	if _gormInstance != nil {
		return _gormInstance.DB, nil
	}

	db, err := gorm.Open(postgres.Open(c.build()), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	_gormInstance = &gormInstance{
		DB: db,
	}

	return _gormInstance.DB, nil
}
