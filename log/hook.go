package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	svcKey = "service"
	ts     = "timestamp"
)

type ServiceHook struct {
	Name string
}

func NewServiceHook(name string) *ServiceHook {
	return &ServiceHook{
		Name: name,
	}
}

func (hook *ServiceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *ServiceHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}

	if hook.Name != "" {
		entry.Data[svcKey] = hook.Name
	}
	entry.Data[ts] = time.Now().Unix()
	return nil
}

// TODO get meta data from request context
