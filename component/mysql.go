package component

import (
	"fmt"
	"sync"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlOnce = &sync.Once{}

type MySqlConfig struct {
	Host                          string `env:""`
	User                          string `env:""`
	Password                      string `env:""`
	DbName                        string `env:""`
	MaxIdleConn                   int    `env:""`
	MaxOpenConn                   int    `env:""`
	DefaultStringSize             uint   `env:""`
	DefaultDatetimePrecision      int    `env:""`
	DisableDatetimePrecision      bool   `env:""`
	DontSupportRenameIndex        bool   `env:""`
	DontSupportRenameColumn       bool   `env:""`
	DontSupportForShareClause     bool   `env:""`
	DontSupportNullAsDefaultValue bool   `env:""`
	DontSupportRenameColumnUnique bool   `env:""`
	SkipInitializeWithVersion     bool   `env:""`
}

type Mysql struct {
	MySqlConfig

	db *gorm.DB
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
		DefaultStringSize: mysql.DefaultStringSize,
		// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DisableDatetimePrecision: mysql.DisableDatetimePrecision,
		// 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameIndex: mysql.DontSupportRenameIndex,
		// 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		DontSupportRenameColumn: mysql.DontSupportRenameColumn,
		// 根据版本自动配置
		SkipInitializeWithVersion: mysql.SkipInitializeWithVersion,
	}

	return mysqlConfig
}
