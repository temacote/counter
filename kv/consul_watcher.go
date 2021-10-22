package kv

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type (
	ConsulWatcher interface {
		Watch(func([]byte) error)
		Get(func([]byte) error) error
	}

	consulWatcher struct {
		conf            *viper.Viper
		logger          *zap.Logger
		client          *api.Client
		lastModifyIndex uint64
		consulKey       string
	}
)

func NewConsulWatcher(client *api.Client, logger *zap.Logger, conf *viper.Viper) ConsulWatcher {
	return &consulWatcher{
		conf:   conf,
		logger: logger,
		client: client,
		consulKey: fmt.Sprintf(
			"/%s/%s",
			conf.GetString("namespace"),
			conf.GetString("service"),
		),
	}
}

func (k *consulWatcher) Watch(handle func(val []byte) error) {
	var (
		logger = k.logger.With(zap.String("subsystem", "kvwatcher"), zap.String("key", k.consulKey))
		err    error
	)

	_ = k.update(handle)
	go func() {
		var ticker = time.NewTicker(time.Second * 10)
		for range ticker.C {
			if err = k.update(handle); err != nil {
				logger.Error("error watch consul key", zap.Error(err))
			}
		}
	}()
}

func (k *consulWatcher) Get(h func([]byte) error) error {
	return k.update(h)
}

func (k *consulWatcher) update(handle func(val []byte) error) (err error) {
	var pair *api.KVPair

	if pair, _, err = k.client.KV().Get(k.consulKey, nil); err != nil {
		return errors.Wrap(err, "error get consul config")
	}

	if pair == nil {
		return errors.Errorf("empty consul key value: %s", k.consulKey)
	}

	if pair.ModifyIndex == k.lastModifyIndex {
		return nil
	}

	if err = handle(pair.Value); err != nil {
		return errors.Wrap(err, "error handle modify kv pair")
	}

	k.lastModifyIndex = pair.ModifyIndex
	return nil
}
