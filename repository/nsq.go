package repository

import (
	"sync"
	"tinamic/conf"
	"tinamic/pkg/nsq"
)

var (
	nsqConfig *nsq.NsqConfig
	nsqOnce   sync.Once
)

func GetNsqConfigInstance() *nsq.NsqConfig {
	nsqOnce.Do(func() {
		cfg := conf.GetConfigInstance()
		nsqConfig = nsq.NewNsqConfig(
			cfg.GetString("message.nsq.host"),
			cfg.GetInt("message.nsq.port"),
		)
	})
	return nsqConfig
}
