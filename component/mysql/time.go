package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type NormalTime struct {
	time.Time
}

func (t NormalTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t *NormalTime) Scan(value interface{}) error {
	if v, ok := value.(time.Time); ok {
		*t = NormalTime{Time: v}
		return nil
	}
	return fmt.Errorf("failed to scan time value: %v", value)
}

func (t NormalTime) Value() (driver.Value, error) {
	return t.Time, nil
}

type DeletedTime struct {
	gorm.DeletedAt
}

func (delete DeletedTime) MarshalJSON() ([]byte, error) {
	if delete.Valid {
		return json.Marshal(delete.Time)
	}
	return json.Marshal("")
}
