package database

import (
	"time"

	"github.com/lie-flat-planet/service-init-tool/component/option"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

type DeletedTime struct {
	gorm.DeletedAt
}

func (delete DeletedTime) MarshalJSON() ([]byte, error) {
	if delete.Valid {
		return json.Marshal(delete.Time)
	}
	return json.Marshal("")
}

type TimestampAt struct {
	CreatedAt   int64  `json:"createdAt" gorm:"autoCreateTime:true;not null"` // 自动创建时间戳(秒)
	UpdatedAt   int64  `json:"updatedAt" gorm:"autoUpdateTime:true;not null"` // 自动更新时间戳(秒)
	CreatedTime string `json:"createdTime" gorm:"-"`                          // 创建时间
	UpdatedTime string `json:"updatedTime" gorm:"-"`                          // 更新时间
	// DeletedAt soft_delete.DeletedAt `json:"deletedAt" gorm:"uniqueIndex:name"`
}

func listOptions(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) []gorm.Option {
	var options = []gorm.Option{
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建数据库外键约束
		},
	}

	for _, ops := range opts {
		options = append(options, ops.(gorm.Option))
	}

	return options

}

func dbLogger() logger.Interface {
	return logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  logger.Info,
		},
	)
}
