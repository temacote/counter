package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"

	"sber_cloud/tw/definition/config"
)

type (
	Redis interface {
		SaveValue(value string, key string, expiration time.Duration) error
		GetAllValues() (data []*storedData, err error)
		ParseAndSave(data []byte, expiration time.Duration) (err error)
	}

	rds struct {
		conf  config.Config
		redis *redis.Client
	}

	storedData struct {
		Key  string `json:"key"`
		Data string `json:"data"`
	}
)

func (d *storedData) GetJson() (data []byte) {
	data, _ = json.Marshal(d)
	return
}

func NewRedis(conf config.Config) Redis {
	return &rds{
		conf: conf,
		redis: redis.NewClient(&redis.Options{
			Addr: conf.GetString("redis.url"),
		}),
	}
}

// SaveValue сохранение значения в редис
func (r *rds) SaveValue(value string, key string, expiration time.Duration) (err error) {
	err = r.redis.Set(key, value, expiration).Err()
	return
}

// GetAllValues получение всех значений из редиса с ключами
func (r *rds) GetAllValues() (data []*storedData, err error) {
	iter := r.redis.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		data = append(data, &storedData{
			Key:  iter.Val(),
			Data: r.redis.Get(iter.Val()).Val(),
		})
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// ParseAndSave парсинг структуры для сохранения в базу и ее непосредственное сохранение
func (r *rds) ParseAndSave(data []byte, expiration time.Duration) (err error) {
	// проверяем приводится ли value к структуре storedData
	var strd *storedData
	err = json.Unmarshal(data, &strd)
	if err != nil {
		return
	}
	err = r.redis.Set(strd.Key, strd.Data, expiration).Err()
	return
}
