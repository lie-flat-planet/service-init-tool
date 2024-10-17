package redis

import (
	"context"
	"fmt"
	"github.com/lie-flat-planet/service-init-tool/component/option"
	"sync"
	"time"

	goRedis "github.com/go-redis/redis/v8"
)

var (
	redisOnce = &sync.Once{}
)

type Config struct {
	Host     string `env:""`
	Username string `env:""`
	Password string `env:""`
	DB       int    `env:""`
	PoolSize int    `env:""`
	Timeout  int    `env:""`
}

type Redis struct {
	Config

	client *goRedis.Client `skip:""`
}

// Init 会被工具自动执行。研发不应该调用该方法
func (r *Redis) Init() error {
	var err error
	redisOnce.Do(
		func() {
			err = r.dial()
		},
	)
	if err != nil {
		return fmt.Errorf("redis init error. %w", err)
	}

	return r.ping()
}

// NewInstance 如果你对实例需要进行新的配置，你可以使用该方法覆写 redis.client
func (r *Redis) NewInstance(opts ...option.ClientOptionInterface[*goRedis.Options, *goRedis.Client]) error {
	if err := r.dial(opts...); err != nil {
		return err
	}

	return r.ping()
}

func (r *Redis) GetClient() *goRedis.Client {
	return r.client
}

func (r *Redis) dial(opts ...option.ClientOptionInterface[*goRedis.Options, *goRedis.Client]) (err error) {
	opt := &goRedis.Options{
		Addr:         r.Host,
		Username:     r.Username,
		Password:     r.Password,
		DB:           r.DB,
		PoolSize:     r.PoolSize,
		ReadTimeout:  time.Second * time.Duration(r.Timeout),
		WriteTimeout: time.Second * time.Duration(r.Timeout),
	}

	var cli *goRedis.Client
	for _, o := range opts {
		if o != nil {
			if err = o.Apply(opt); err != nil {
				return err
			}
			defer func() {
				if afterErr := o.AfterInitialize(cli); afterErr != nil {
					err = afterErr
				}
			}()
		}
	}

	cli = goRedis.NewClient(opt)
	r.client = cli

	return
}

func (r *Redis) ping() error {
	ctx := context.Background()

	err := r.client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis ping failed. err:%v", err)
	}
	return nil
}
