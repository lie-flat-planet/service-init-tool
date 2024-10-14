package mysql

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t *Time) Scan(value interface{}) error {
	if v, ok := value.(time.Time); ok {
		*t = Time{Time: v}
		return nil
	}
	return fmt.Errorf("failed to scan time value: %v", value)
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}
