# Server-side component for Clockwork browser extension written in GO

## Attach and use data-sources to clockwork instance.

```go
package main
import (

"github.com/anton-shumanski/clockwork"
"github.com/go-redis/redis"
)


func main()  {
 client := redis.NewClient(&redis.Options{
 		Addr:     "127.0.0.1:6379", //your address
 		Password: "", // your password
 		DB:       0, 
 	})
 	
 profiler := clockwork.Clockwork{RedisStorageProvider: client}
}
```

Mysql data source
```go
var mysqlDataSource dataSource.QueryLoggerDataSourceInterface = &dataSource.MysqlDataSource{}
profiler.AddDataSource(mysqlDataSource)

mysqlDataSource.LogQuery("Select * from users where id = ?", 12.224, []string{"a"})
mysqlDataSource.LogQuery("Select * from address where id = ?", 1, []string{"a"})
```

Redis data source
```go
var redisDataSource dataSource.CommandLoggerDataSourceInterface = &dataSource.RedisDataSource{}
profiler.AddDataSource(redisDataSource)

redisDataSource.LogCommand("hSet", "test_key_1", 0.12)
redisDataSource.LogCommand("hGet", "test_key_2", 0.15)
```

Cache data source
```go
var cacheDataSource dataSource.CacheLoggerDataSourceInterface = &dataSource.CacheDataSource{}
profiler.AddDataSource(cacheDataSource)

cacheDataSource.LogCache("miss", "price", "30.10$", 12.22, 3000)
```

Timeline data source
```go
var timelineDataSource dataSource.TimelineLoggerDataSourceInterface = &dataSource.TimelineDataSource{}
profiler.SetTimeLineDataSource(timelineDataSource)

profiler.GetTimeLineDataSource().StartEvent("Event_11", "My first event desc")
//put some logic here
profiler.GetTimeLineDataSource().EndEvent("Event_11")

profiler.GetTimeLineDataSource().StartEvent("Event_22", "My second event desc")
//put some logic here
profiler.GetTimeLineDataSource().EndEvent("Event_22")
```

Request data source
```go
var requestDataSource dataSource.RequestLoggerDataSourceInterface = &dataSource.RequestResponseDataSource{}
profiler.SetRequestDataSource(requestDataSource)

//in the begining of request
profiler.GetRequestDataSource().SetStartTime(time.Now())
profiler.GetRequestDataSource().StartMemoryUsage()
profiler.GetRequestDataSource().SetController("HomeController", "IndexAction")

//at the end of request	
profiler.GetRequestDataSource().SetResponseTime(time.Now())
profiler.GetRequestDataSource().SetResponseStatus(200)
```

Logger(debugger) data source
```go
var loggerDataSource dataSource.LoggerDataSourceInterface = &dataSource.LoggerDataSource{}
profiler.SetLoggerDataSource(loggerDataSource)

profiler.GetLoggerDataSource().LogDebugString("test payment", "payment method works")
test1 := make([]string, 2)
test1[0] = "item1"
test1[1] = "item 2"
profiler.GetLoggerDataSource().LogDebugSlice("cart", test1)

users := make(map[string]string)
users["Adam"] = "33 ages"
users["Fred"] = "15 ages"
profiler.GetLoggerDataSource().LogDebugMap("users", users)	
```

## Before end of the request
You should call those 2 methods
```go
profiler.Resolve()
profiler.SaveData()
```

The last thing u need to do is to send 2 special headers into the response:
```go
c.Writer.Header().Set("X-Clockwork-Id", profiler.GetUniqueId())
c.Writer.Header().Set("X-Clockwork-Version", "4.0.13")
```