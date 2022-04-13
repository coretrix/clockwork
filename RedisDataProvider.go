package clockwork

import (
	"encoding/json"
	dataSource "github.com/coretrix/clockwork/data-source"
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisDataProvider struct {
	RedisStorageProvider *redis.Client
}

func (provider *RedisDataProvider) Get(key string, id string) dataSource.DataBuffer {
	result, err := provider.RedisStorageProvider.HGet(key, id).Result()
	if err != nil {
		panic(err)
	}

	provider.RedisStorageProvider.Expire(key, time.Minute*5)
	var raw dataSource.DataBuffer
	err = json.Unmarshal([]byte(result), &raw)

	return raw
}

func (provider *RedisDataProvider) Set(key string, id string, data *dataSource.DataBuffer) {
	jsonString, _ := json.Marshal(data)
	err := provider.RedisStorageProvider.HSet(key, id, jsonString).Err()
	if err != nil {
		panic(err)
	}
}

type DataProviderInterface interface {
	Get(key string, id string) dataSource.DataBuffer
	Set(key string, id string, data *dataSource.DataBuffer)
}
