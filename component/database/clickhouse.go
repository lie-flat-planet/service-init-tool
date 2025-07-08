package database

import (
	"fmt"
	"sync"

	"github.com/lie-flat-planet/service-init-tool/component/option"

	gormClickhouse "gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var (
	clickhouseOnce = &sync.Once{}
)

type ClickhouseConf struct {
	Host        string `env:""`
	User        string `env:""`
	Password    string `env:""`
	DbName      string `env:""`
	Port        int    `env:""`
	MaxIdleConn int    `env:""`
	MaxOpenConn int    `env:""`
	IgnoreLog   bool   `env:""`
	// models is the gorm Models
	models []any
}

type Clickhouse struct {
	ClickhouseConf

	db *gorm.DB `skip:""`
}

// Init 会被工具自动执行。研发不应该调用该方法
func (clickhouse *Clickhouse) Init() error {
	var err error
	clickhouseOnce.Do(
		func() {
			err = clickhouse.dialAndSetConn()
		},
	)
	if err != nil {
		return fmt.Errorf("clickhouse init error. %w", err)
	}

	return clickhouse.ping()
}

func (clickhouse *Clickhouse) GetDB() *gorm.DB {
	return clickhouse.db
}

// NewInstance 如果你对实例需要进行新的配置，你可以使用该方法覆写 clickhouse.db
func (clickhouse *Clickhouse) NewInstance(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	if err := clickhouse.dialAndSetConn(opts...); err != nil {
		return err
	}

	return clickhouse.ping()
}

func (clickhouse *Clickhouse) NewSession(cfg ...*gorm.Session) *gorm.DB {
	if len(cfg) < 1 {
		return clickhouse.db.Session(&gorm.Session{})
	}

	return clickhouse.db.Session(cfg[0])
}

func (clickhouse *Clickhouse) AppendModel(model ...any) {
	clickhouse.models = append(clickhouse.models, model...)
}

// MigrateAll 配合 AppendModel, 一般用与迁移所有表。
// 如果一次性只是迁移部分表，建议使用 MigrateTable
func (clickhouse *Clickhouse) MigrateAll() error {
	return clickhouse.db.AutoMigrate(clickhouse.models...)
}

// MigrateTable 用于迁移部分表
func (clickhouse *Clickhouse) MigrateTable(model ...any) error {
	return clickhouse.db.AutoMigrate(model...)
}

func (clickhouse *Clickhouse) dialAndSetConn(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	var gormOptions = listOptions(opts...)

	if !clickhouse.IgnoreLog {
		gormOptions = append(gormOptions, &gorm.Config{
			Logger: dbLogger(),
		})
	}

	// 初始化 clickhouse client
	cfg := clickhouse.formConfig()
	db, err := gorm.Open(gormClickhouse.New(cfg), gormOptions...)
	if err != nil {
		return fmt.Errorf("clickhouse init failed. err:%w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(clickhouse.MaxIdleConn)
	sqlDB.SetMaxOpenConns(clickhouse.MaxOpenConn)

	clickhouse.db = db

	return nil
}

func (clickhouse *Clickhouse) formConfig() gormClickhouse.Config {
	// "clickhouse://gorm:gorm@localhost:9942/gorm?dial_timeout=10s&read_timeout=20s"
	if clickhouse.Port == 0 {
		clickhouse.Port = 9000
	}
	dsn := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%d/%s?dial_timeout=10s&max_execution_time=60",
		clickhouse.User,
		clickhouse.Password,
		clickhouse.Host,
		clickhouse.Port,
		clickhouse.DbName,
	)

	clickhouseConfig := gormClickhouse.Config{
		// DSN data source name
		DSN: dsn,
	}

	return clickhouseConfig
}

func (clickhouse *Clickhouse) ping() error {
	sqlDB, err := clickhouse.db.DB()
	if err != nil {
		return fmt.Errorf("get sqldb error. %w", err)
	}

	return sqlDB.Ping()
}
