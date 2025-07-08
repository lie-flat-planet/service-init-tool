package database

import (
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

var clickhouse = &Clickhouse{
	ClickhouseConf: ClickhouseConf{
		Host:     "localhost",
		Port:     9001,
		User:     "aaa",
		Password: "xxx",
		DbName:   "app",
	},
}

func TestClickhouse_GetDB(t *testing.T) {
	err := clickhouse.Init()
	assert.NoError(t, err)

	db := clickhouse.GetDB()
	assert.NotNil(t, db)
}

// Scalar 模型定义
type Scalar struct {
	UID          uint64  `gorm:"column:uid;type:UInt64;not null"`
	ProjectID    string  `gorm:"column:projectId;type:String;not null"`
	ExperimentID string  `gorm:"column:experimentId;type:String;not null"`
	Key          string  `gorm:"column:key;type:String;not null"`
	Epoch        uint64  `gorm:"column:epoch;type:UInt64;not null"`
	Step         uint64  `gorm:"column:step;type:UInt64;not null"`
	Value        float64 `gorm:"column:value;type:Float64;not null"`
	Timestamp    Time    `gorm:"column:timestamp;type:DateTime64(3);not null"`
	CreatedAt    Time    `gorm:"column:createdAt;type:DateTime;default:now();not null"`
}

// TableName 指定表名
func (Scalar) TableName() string {
	return "scalar"
}

// GetScalarByExperimentID 根据实验ID查询标量数据
func GetScalarByExperimentID(experimentID string, db *gorm.DB) ([]Scalar, error) {
	var scalars []Scalar
	err := db.Where("experimentId = ?", experimentID).Find(&scalars).Error
	if err != nil {
		return nil, err
	}
	return scalars, nil
}

func TestClickhouse_GetScalarByExperimentID(t *testing.T) {
	err := clickhouse.Init()
	assert.NoError(t, err)

	db := clickhouse.GetDB()
	assert.NotNil(t, db)

	scalars, err := GetScalarByExperimentID("j5q7eo8u9tfye4e64uiaf", db)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, scalars)

}
