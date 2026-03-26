package user_migrate

import "gorm.io/gorm"

func Migrate(db *gorm.DB, models []any) error {

	if len(models) == 0 {
		return nil
	}

	return db.AutoMigrate(models...)
}