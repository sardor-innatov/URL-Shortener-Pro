package helper

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JsonObject map[string]any

func (j JsonObject) Value() (driver.Value, error) {
	return json.Marshal(j)

}
func (j *JsonObject) Scan(value any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}
	return json.Unmarshal(bytes, j)
}

func (JsonObject) GormDataType() string {
	return "json"
}

func (JsonObject) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "jsonb"
}
