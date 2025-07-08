package database

import (
	"fmt"
	"sync"

	"github.com/lie-flat-planet/service-init-tool/component/option"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlOnce = &sync.Once{}
)

type MysqlConf struct {
	Host        string `env:""`
	User        string `env:""`
	Password    string `env:""`
	DbName      string `env:""`
	MaxIdleConn int    `env:""`
	MaxOpenConn int    `env:""`
	IgnoreLog   bool   `env:""`
	// models is the gorm Models
	models []any
}

type Mysql struct {
	MysqlConf

	db *gorm.DB `skip:""`
}

// Init 会被工具自动执行。研发不应该调用该方法
func (mysql *Mysql) Init() error {
	var err error
	mysqlOnce.Do(
		func() {
			err = mysql.dialAndSetConn()
		},
	)
	if err != nil {
		return fmt.Errorf("mysql init error. %w", err)
	}

	return mysql.ping()
}

func (mysql *Mysql) GetDB() *gorm.DB {
	return mysql.db
}

// NewInstance 如果你对实例需要进行新的配置，你可以使用该方法覆写 mysql.db
func (mysql *Mysql) NewInstance(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	if err := mysql.dialAndSetConn(opts...); err != nil {
		return err
	}

	return mysql.ping()
}

func (mysql *Mysql) NewSession(cfg ...*gorm.Session) *gorm.DB {
	if len(cfg) < 1 {
		return mysql.db.Session(&gorm.Session{})
	}

	return mysql.db.Session(cfg[0])
}

func (mysql *Mysql) AppendModel(model ...any) {
	mysql.models = append(mysql.models, model...)
}

// MigrateAll 配合 AppendModel, 一般用与迁移所有表。
// 如果一次性只是迁移部分表，建议使用 MigrateTable
func (mysql *Mysql) MigrateAll() error {
	return mysql.db.AutoMigrate(mysql.models...)
}

// MigrateTable 用于迁移部分表
func (mysql *Mysql) MigrateTable(model ...any) error {
	return mysql.db.AutoMigrate(model...)
}

func (mysql *Mysql) dialAndSetConn(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	var gormOptions = listOptions(opts...)

	if !mysql.IgnoreLog {
		gormOptions = append(gormOptions, &gorm.Config{
			Logger: dbLogger(),
		})
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

func (mysql *Mysql) ping() error {
	sqlDB, err := mysql.db.DB()
	if err != nil {
		return fmt.Errorf("get sqldb error. %w", err)
	}

	return sqlDB.Ping()
}
