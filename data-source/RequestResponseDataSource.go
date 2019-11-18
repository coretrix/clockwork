package dataSource

import (
	"runtime"
	"time"
)

type RequestLoggerDataSourceInterface interface {
	DataSource
	RequestLoggerInterface
}

type RequestLoggerInterface interface {
	SetStartTime(start time.Time)
	SetResponseTime(responseTime time.Time)
	SetResponseStatus(status int16)
	StartMemoryUsage()
	EndMemoryUsage()
	SetController(controller string, method string)
}

type RequestResponseDataSource struct {
	commands map[string]interface{}
	startTime float64
	startMemoryUsage uint64
	endMemoryUsage uint64
	responseTime float64
	responseStatus int16
	description string
	controller string
	data []string
}

func (source *RequestResponseDataSource) SetController(controller string, method string)  {
	source.controller = controller+ "@" +method
}

func (source *RequestResponseDataSource) SetStartTime(start time.Time)  {
	source.startTime = MicroTime(start)
}

func (source *RequestResponseDataSource) SetResponseTime(responseTime time.Time)  {
	source.responseTime = MicroTime(responseTime)
}

func (source *RequestResponseDataSource) StartMemoryUsage()  {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	source.startMemoryUsage = m.Alloc
}

func (source *RequestResponseDataSource) EndMemoryUsage()  {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	source.endMemoryUsage = m.Alloc
}

func (source *RequestResponseDataSource) SetResponseStatus(responseStatus int16)  {
	source.responseStatus = responseStatus
}

func (source *RequestResponseDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.Time = source.startTime
	dataBuffer.Controller = source.controller
	dataBuffer.MemoryUsage = source.endMemoryUsage - source.startMemoryUsage
	dataBuffer.ResponseStatus = source.responseStatus
	dataBuffer.ResponseTime = source.responseTime
	dataBuffer.ResponseDuration = (	source.responseTime  * 1000 ) - (source.startTime * 1000)
}
