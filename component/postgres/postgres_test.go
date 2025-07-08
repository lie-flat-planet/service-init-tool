package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Key 模型定义
type Key struct {
	ID        int    `gorm:"primaryKey;column:id"`
	Name      string `gorm:"column:name;not null"`
	Key       string `gorm:"column:key;not null"`
	CreatedAt string `gorm:"column:createdAt;not null;default:CURRENT_TIMESTAMP"`
	UserID    int    `gorm:"column:userId;not null"`
}

// TableName 指定表名
func (Key) TableName() string {
	return "Key"
}

var pg = &Postgres{
	Config: Config{
		Host:     "localhost",
		User:     "swanlab",
		Password: "xxx",
		DbName:   "app",
		Port:     5432,
	},
}

func TestPostgres_GetDB(t *testing.T) {
	err := pg.Init()
	assert.NoError(t, err)

	db := pg.GetDB()
	assert.NotNil(t, db)
}

// GetKeyByID 根据ID查询单个密钥
func KeyGetKeyByID(id int, db *gorm.DB) (*Key, error) {
	var key Key
	err := db.Where("id = ?", id).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func TestPostgres_GetKeyByID(t *testing.T) {
	err := pg.Init()
	assert.NoError(t, err)

	db := pg.GetDB()
	assert.NotNil(t, db)

	k, err := KeyGetKeyByID(1, db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(k)
}
