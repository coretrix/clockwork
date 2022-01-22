package clockwork

import (
	"encoding/json"
	"fmt"
	dataSource "github.com/anton-shumanski/clockwork/data-source"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const key = "profiler_store"

type Clockwork struct {
	DataProvider       DataProviderInterface
	dataSources        []dataSource.DataSource
	id                 string
	timeLineDataSource *dataSource.TimelineLoggerDataSourceInterface
	requestDataSource  *dataSource.RequestLoggerDataSourceInterface
	loggerDataSource   *dataSource.LoggerDataSourceInterface
	databaseDataSource *dataSource.QueryLoggerDataSourceInterface
	data               *dataSource.DataBuffer
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

	a, err := json.Marshal(clockwork.data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(a))
	return clockwork.data
}

func (clockwork *Clockwork) SaveData() {
	clockwork.DataProvider.Set(key, clockwork.GetUniqueId(), clockwork.Resolve())
}

func (clockwork *Clockwork) GetSavedData(id string) dataSource.DataBuffer {
	return clockwork.DataProvider.Get(key, id)
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

func (clockwork *Clockwork) GetDatabaseDataSource() dataSource.QueryLoggerInterface {
	return *clockwork.databaseDataSource
}

func (clockwork *Clockwork) SetDatabaseDataSource(source dataSource.QueryLoggerDataSourceInterface) {
	clockwork.databaseDataSource = &source
	clockwork.AddDataSource(source)
}

func (clockwork *Clockwork) GetLoggerDataSource() dataSource.LoggerDataSourceInterface {
	return *clockwork.loggerDataSource
}
