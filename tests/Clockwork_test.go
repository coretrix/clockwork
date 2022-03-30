package tests

import (
	"fmt"
	"github.com/anton-shumanski/clockwork"
	dataSource "github.com/anton-shumanski/clockwork/data-source"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClockwork_GetData(t *testing.T) {
	profiler := clockwork.Clockwork{}
	var mysqlDataSource dataSource.QueryLoggerDataSourceInterface = new(dataSource.DatabaseDataSource)
	var redisDataSource dataSource.CommandLoggerDataSourceInterface = new(dataSource.RedisDataSource)
	var cacheDataSource dataSource.CacheLoggerDataSourceInterface = new(dataSource.CacheDataSource)
	var customDataSource dataSource.UserDataSourceInterface = new(dataSource.UserDataDataSource)
	var timelineDataSource dataSource.TimelineLoggerDataSourceInterface = new(dataSource.TimelineDataSource)
	var requestResponseDataSource dataSource.RequestLoggerDataSourceInterface = new(dataSource.RequestResponseDataSource)

	profiler.AddDataSource(redisDataSource)
	profiler.AddDataSource(cacheDataSource)
	profiler.AddDataSource(customDataSource)
	profiler.SetTimeLineDataSource(timelineDataSource)
	profiler.SetRequestDataSource(requestResponseDataSource)
	profiler.SetDatabaseDataSource(mysqlDataSource)

	requestResponseDataSource.SetStartTime(time.Now())
	var bind1 []interface{}
	var bind2 []interface{}
	middleware := []string{"Authorize", "Normalization", "Guard", "Handler"}

	bind2 = append(bind2, 1, 2, "test param")
	mysqlDataSource.LogQuery("mysql", "SELECT * FROM users", 12.224, bind1)
	mysqlDataSource.LogQuery("mysql", "SELECT * FROM address where id = ?", 1, bind2)

	redisDataSource.LogCommand("hSet", "test_key_1", 0.12)
	redisDataSource.LogCommand("hGet", "test_key_1", 0.15)

	cacheDataSource.LogCacheMiss("pool-1", "hGet", "test_key", "30cm", 1, 12.22, 3000)

	timelineDataSource.StartEvent("Request_11", "My first request desc")
	timelineDataSource.EndEvent("Request_11")

	timelineDataSource.StartEvent("Request_22", "My second request desc")
	timelineDataSource.EndEvent("Request_22")

	requestResponseDataSource.SetResponseTime(time.Now())
	requestResponseDataSource.SetResponseStatus(200)
	requestResponseDataSource.SetMiddleware(middleware)

	customDataSource.SetShowAs("table")
	customDataSource.SetTitle("test")

	response := profiler.Resolve()
	fmt.Println(response.UserData)
	assert.Equal(t, len(response.TimelineData), 2)
	assert.Equal(t, len(response.DatabaseQueries), 2)
	assert.Equal(t, len(response.CacheQueries), 1)
	assert.Equal(t, len(response.RedisCommands), 2)
	assert.Equal(t, response.DatabaseQueriesCount, 2)
	assert.Equal(t, response.ResponseStatus, int16(200))
	assert.Equal(t, response.CacheReads, int16(1))
	assert.Equal(t, response.CacheHits, int16(0))
	assert.Equal(t, response.CacheDeletes, int16(0))
	assert.Equal(t, response.CacheWrites, int16(0))
	assert.Equal(t, response.Middleware, middleware)
}
