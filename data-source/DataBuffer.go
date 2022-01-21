package dataSource

type DataBuffer struct {
	Log                  []interface{}          `json:"log"`
	DatabaseQueries      []interface{}          `json:"databaseQueries"`
	DatabaseQueriesCount int                    `json:"databaseQueriesCount"`
	DatabaseDuration     float32                `json:"databaseDuration"`
	RedisCommands        []interface{}          `json:"redisCommands"`
	CacheQueries         []interface{}          `json:"cacheQueries"`
	CacheTime            float32                `json:"cacheTime"`
	CacheDeletes         int16                  `json:"cacheDeletes"`
	CacheWrites          int16                  `json:"cacheWrites"`
	CacheReads           int16                  `json:"cacheReads"`
	CacheHits            int16                  `json:"cacheHits"`
	TimelineData         map[string]interface{} `json:"timelineData"`
	Time                 float64                `json:"time"`
	Controller           string                 `json:"controller"`
	Middleware           []string               `json:"middleware"`
	MemoryUsage          uint64                 `json:"memoryUsage"`
	ResponseTime         float64                `json:"responseTime"`
	ResponseDuration     float64                `json:"responseDuration"`
	ResponseStatus       int16                  `json:"responseStatus"`
	UserData             []*UserDataDataSource  `json:"userData"`
}
