package database

import (
	"fmt"
	"sync"

	"github.com/lie-flat-planet/service-init-tool/component/option"

	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresOnce = &sync.Once{}
)

type PostgresConf struct {
	Host                 string `env:""`
	User                 string `env:""`
	Password             string `env:""`
	DbName               string `env:""`
	Port                 int    `env:""`
	SSLMode              string `env:""`
	MaxIdleConn          int    `env:""`
	MaxOpenConn          int    `env:""`
	IgnoreLog            bool   `env:""`
	PreferSimpleProtocol bool   `env:""`
	// models is the gorm Models
	models []any
}

type Postgres struct {
	PostgresConf

	db *gorm.DB `skip:""`
}

// Init 会被工具自动执行。研发不应该调用该方法
func (postgres *Postgres) Init() error {
	var err error
	postgresOnce.Do(
		func() {
			err = postgres.dialAndSetConn()
		},
	)
	if err != nil {
		return fmt.Errorf("postgres init error. %w", err)
	}

	return postgres.ping()
}

func (postgres *Postgres) GetDB() *gorm.DB {
	return postgres.db
}

// NewInstance 如果你对实例需要进行新的配置，你可以使用该方法覆写 postgres.db
func (postgres *Postgres) NewInstance(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	if err := postgres.dialAndSetConn(opts...); err != nil {
		return err
	}

	return postgres.ping()
}

func (postgres *Postgres) NewSession(cfg ...*gorm.Session) *gorm.DB {
	if len(cfg) < 1 {
		return postgres.db.Session(&gorm.Session{})
	}

	return postgres.db.Session(cfg[0])
}

func (postgres *Postgres) AppendModel(model ...any) {
	postgres.models = append(postgres.models, model...)
}

// MigrateAll 配合 AppendModel, 一般用与迁移所有表。
// 如果一次性只是迁移部分表，建议使用 MigrateTable
func (postgres *Postgres) MigrateAll() error {
	return postgres.db.AutoMigrate(postgres.models...)
}

// MigrateTable 用于迁移部分表
func (postgres *Postgres) MigrateTable(model ...any) error {
	return postgres.db.AutoMigrate(model...)
}

func (postgres *Postgres) dialAndSetConn(opts ...option.ClientOptionInterface[*gorm.Config, *gorm.DB]) error {
	var gormOptions = listOptions(opts...)

	if !postgres.IgnoreLog {
		gormOptions = append(gormOptions, &gorm.Config{
			Logger: dbLogger(),
		})
	}

	// 初始化 postgres client
	cfg := postgres.formConfig()
	db, err := gorm.Open(gormPostgres.New(cfg), gormOptions...)
	if err != nil {
		return fmt.Errorf("postgres init failed. err:%w", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(postgres.MaxIdleConn)
	sqlDB.SetMaxOpenConns(postgres.MaxOpenConn)

	postgres.db = db

	return nil
}

func (postgres *Postgres) formConfig() gormPostgres.Config {
	if postgres.SSLMode == "" {
		postgres.SSLMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		postgres.Host,
		postgres.User,
		postgres.Password,
		postgres.DbName,
		postgres.Port,
		postgres.SSLMode,
	)

	postgresConfig := gormPostgres.Config{
		// DSN data source name
		DSN: dsn,
		// 禁用预编译语句
		PreferSimpleProtocol: postgres.PreferSimpleProtocol,
	}

	return postgresConfig
}

func (postgres *Postgres) ping() error {
	sqlDB, err := postgres.db.DB()
	if err != nil {
		return fmt.Errorf("get sqldb error. %w", err)
	}

	return sqlDB.Ping()
}
