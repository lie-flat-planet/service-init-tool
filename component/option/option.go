package option

import (
	goRedis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ClientOptionInterface 主要目的将第三方库的配置给暴露出来
// TODO: 后续还有其组件类型可以追加
type ClientOptionInterface[T any | *gorm.Config | *goRedis.Options, AT any | *gorm.DB | *goRedis.Client] interface {
	Apply(T) error
	AfterInitialize(AT) error
}
