package component

import (
	"fmt"
	"sync"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlOnce = &sync.Once{}

type MySqlConfig struct {
	Host        string `env:""`
	User        string `env:""`
	Password    string `env:""`
	DbName      string `env:""`
	MaxIdleConn int    `env:""`
	MaxOpenConn int    `env:""`
}

type Mysql struct {
	MySqlConfig

	db *gorm.DB `skipEnv:""`
}

func (mysql *Mysql) GetDB() *gorm.DB {
	return mysql.db
}

func (mysql *Mysql) NewInstance(opts ...ClientOptionInterface[*gorm.Config, *gorm.DB]) (*gorm.DB, error) {
	var err error
	mysqlOnce.Do(
		func() {
			err = mysql.dialAndSetConn(opts...)
		},
	)

	return mysql.db, err
}

func (mysql *Mysql) dialAndSetConn(opts ...ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	var gormOptions []gorm.Option
	for _, ops := range opts {
		gormOptions = append(gormOptions, ops.(gorm.Option))
	}

	// 初始化 mysql client
	cfg := mysql.formConfig()
	db, err := gorm.Open(gormMysql.New(cfg), gormOptions...)
	if err != nil {
		return fmt.Errorf("mysql init failed. err:%w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(mysql.MaxIdleConn)
	sqlDB.SetMaxOpenConns(mysql.MaxOpenConn)

	mysql.db = db

	return nil
}

func (mysql *Mysql) formConfig() gormMysql.Config {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysql.User,
		mysql.Password,
		mysql.Host,
		mysql.DbName,
	)

	mysqlConfig := gormMysql.Config{
		// DSN data source name
		DSN: dsn,
		// string 类型字段的默认长度
		DefaultStringSize: 256,
		// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DisableDatetimePrecision: true,
	}

	return mysqlConfig
}
