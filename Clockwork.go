package clockwork

import (
	"encoding/json"
	dataSource "github.com/anton-shumanski/clockwork/data-source"
	"github.com/go-redis/redis/v7"
	"math/rand"
	"strconv"
	"time"
)

const key = "profiler_store"

type Clockwork struct {
	RedisStorageProvider *redis.Client
	dataSources          []dataSource.DataSource
	id                   string
	timeLineDataSource   *dataSource.TimelineLoggerDataSourceInterface
	requestDataSource    *dataSource.RequestLoggerDataSourceInterface
	loggerDataSource     *dataSource.LoggerDataSourceInterface
	data                 *dataSource.DataBuffer
}

func (clockwork *Clockwork) AddDataSource(source dataSource.DataSource) *dataSource.DataSource {
	clockwork.dataSources = append(clockwork.dataSources, source)
	return &source
}

func (clockwork *Clockwork) Resolve() *dataSource.DataBuffer {
	clockwork.data = &dataSource.DataBuffer{}
	for _, source := range clockwork.dataSources {
		source.Resolve(clockwork.data)
	}

	return clockwork.data
}

func (clockwork *Clockwork) SaveData() {
	jsonString, _ := json.Marshal(clockwork.data)
	err := clockwork.RedisStorageProvider.HSet(key, clockwork.GetUniqueId(), jsonString).Err()
	if err != nil {
		panic(err)
	}
}

func (clockwork *Clockwork) GetSavedData(id string) dataSource.DataBuffer {
	result, err := clockwork.RedisStorageProvider.HGet(key, id).Result()
	if err != nil {
		panic(err)
	}

	clockwork.RedisStorageProvider.Expire(key, time.Minute*5)
	var raw dataSource.DataBuffer
	err = json.Unmarshal([]byte(result), &raw)

	return raw
}

func (clockwork *Clockwork) GetUniqueId() string {
	if clockwork.id != "" {
		return clockwork.id
	}

	now := time.Now()
	sec := now.Unix()
	clockwork.id = strconv.FormatInt(sec, 10) + "-" + strconv.FormatInt(rand.Int63(), 10)

	return clockwork.id
}

func (clockwork *Clockwork) SetTimeLineDataSource(source dataSource.TimelineLoggerDataSourceInterface) {
	clockwork.timeLineDataSource = &source
	clockwork.AddDataSource(source)
}

func (clockwork *Clockwork) GetTimeLineDataSource() dataSource.TimelineLoggerDataSourceInterface {
	return *clockwork.timeLineDataSource
}

func (clockwork *Clockwork) SetRequestDataSource(source dataSource.RequestLoggerDataSourceInterface) {
	clockwork.requestDataSource = &source
	clockwork.AddDataSource(source)
}

func (clockwork *Clockwork) GetRequestDataSource() dataSource.RequestLoggerDataSourceInterface {
	return *clockwork.requestDataSource
}

func (clockwork *Clockwork) SetLoggerDataSource(source dataSource.LoggerDataSourceInterface) {
	clockwork.loggerDataSource = &source
	clockwork.AddDataSource(source)
}

func (clockwork *Clockwork) GetLoggerDataSource() dataSource.LoggerDataSourceInterface {
	return *clockwork.loggerDataSource
}
