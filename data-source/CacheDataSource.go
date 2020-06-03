package dataSource

import "fmt"

const (
	CacheHit    = "hit"
	CacheMiss   = "miss"
	CacheWrite  = "write"
	CacheDelete = "delete"
)

type CacheLoggerDataSourceInterface interface {
	DataSource
	CacheLoggerInterface
}

type CacheLoggerInterface interface {
	LogCache(typeParam, action string, key string, value string, duration float32, expiration float32)
}

type cacheDataStructure struct {
	Type       string  `json:"type"`
	Key        string  `json:"key"`
	Value      string  `json:"value"`
	Expiration float32 `json:"expiration"`
	Duration   float32 `json:"duration"`
	Connection string  `json:"connection"`
}

type CacheDataSource struct {
	commands      []interface{}
	totalDuration float32
	cacheReads    int16
	cacheDeletes  int16
	cacheWrites   int16
	cacheHits     int16
}

func (source *CacheDataSource) LogCache(typeParam, action string, key string, value string, duration float32, expiration float32) {
	switch typeParam {
	case CacheHit:
		source.cacheHits += 1
		source.cacheReads += 1
	case CacheMiss:
		source.cacheReads += 1
	case CacheWrite:
		source.cacheWrites += 1
	case CacheDelete:
		source.cacheDeletes += 1
	default:
		panic("There is no supported type " + typeParam)
	}

	structure := cacheDataStructure{
		Type:       fmt.Sprintf("%s %s", action, typeParam),
		Key:        key,
		Value:      value,
		Expiration: expiration,
		Duration:   duration,
		Connection: "test-connection",
	}

	source.totalDuration += duration
	source.commands = append(source.commands, &structure)
}

func (source *CacheDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.CacheQueries = source.commands

	dataBuffer.CacheDeletes = source.cacheDeletes
	dataBuffer.CacheHits = source.cacheHits
	dataBuffer.CacheReads = source.cacheReads
	dataBuffer.CacheWrites = source.cacheWrites
	dataBuffer.CacheTime = source.totalDuration
}
