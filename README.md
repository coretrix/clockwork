# Server-side component for Clockwork browser extension written in GO

## Attach and use data-sources to clockwork instance.

```gotemplate
client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", //your address
		Password: "", // your password
		DB:       0, 
	})
	
profiler := clockwork.Clockwork{RedisStorageProvider: client}
```

Mysql data source
```gotemplate
var mysqlDataSource dataSource.QueryLoggerDataSourceInterface = &dataSource.MysqlDataSource{}
profiler.AddDataSource(mysqlDataSource)

mysqlDataSource.LogQuery("Select * from users where id = ?", 12.224, []string{"a"})
mysqlDataSource.LogQuery("Select * from address where id = ?", 1, []string{"a"})
```

Redis data source
```gotemplate
var redisDataSource dataSource.CommandLoggerDataSourceInterface = &dataSource.RedisDataSource{}
profiler.AddDataSource(redisDataSource)

redisDataSource.LogCommand("hSet", "test_key_1", 0.12)
redisDataSource.LogCommand("hGet", "test_key_2", 0.15)
```

Cache data source
```gotemplate
var cacheDataSource dataSource.CacheLoggerDataSourceInterface = &dataSource.CacheDataSource{}
profiler.AddDataSource(cacheDataSource)

cacheDataSource.LogCache("miss", "price", "30.10$", 12.22, 3000)
```

Timeline data source
```gotemplate
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
```gotemplate
var requestDataSource dataSource.RequestLoggerDataSourceInterface = &dataSource.RequestResponseDataSource{}
profiler.SetRequestDataSource(requestDataSource)

//in the begining of request
profiler.GetRequestDataSource().SetStartTime(time.Now())
profiler.GetRequestDataSource().StartMemoryUsage()
profiler.GetRequestDataSource().SetController("HomeController", "IndexAction")

//at the end of request	
profiler.GetRequestDataSource().SetStartTime(time.Now())
profiler.GetRequestDataSource().StartMemoryUsage()
profiler.GetRequestDataSource().SetController("HomeController", "IndexAction")
```

Logger(debugger) data source
```gotemplate
var loggerDataSource dataSource.LoggerDataSourceInterface = &dataSource.LoggerDataSource{}
profiler.SetLoggerDataSource(loggerDataSource)

profiler.GetLoggerDataSource().LogDebugString("test payment", "payment method works")
test1 := make([]string, 2)
test1[0] = "item1"
test1[1] = "item 2"
profiler.GetLoggerDataSource().LogDebugSlice("cart", test1)

users := make(map[string]string)
slove["Adam"] = "33 ages"
slove["Fred"] = "15 agas"
profiler.GetLoggerDataSource().LogDebugMap("users", users)	
```

## Before end of the request
You should call those 2 methods
```gotemplate
profiler.Resolve()
profiler.SaveData()
```

The last thing u need to do is to send 2 special headers into the response:
```gotemplate
c.Writer.Header().Set("X-Clockwork-Id", profiler.GetUniqueId())
c.Writer.Header().Set("X-Clockwork-Version", "4.0.13")
```