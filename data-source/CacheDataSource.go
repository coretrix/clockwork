package datasource

import "fmt"

const (
	CacheHit    = "hit"
	CacheWrite  = "write"
	CacheDelete = "delete"
)

type CacheLoggerDataSourceInterface interface {
	DataSource
	CacheLoggerInterface
}

type CacheLoggerInterface interface {
	LogCache(connection, typeParam, action string, key string, value string, duration float32, expiration float32)
	LogCacheMiss(connection, action string, key string, value string, misses int, duration float32, expiration float32)
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

func (source *CacheDataSource) LogCache(connection, typeParam, action string, key string, value string, duration float32, expiration float32) {
	switch typeParam {
	case CacheHit:
		source.cacheHits++
		source.cacheReads++
	case CacheWrite:
		source.cacheWrites++
	case CacheDelete:
		source.cacheDeletes++
	default:
		panic("there is no supported type " + typeParam)
	}

	structure := cacheDataStructure{
		Type:       fmt.Sprintf("%s %s", typeParam, action),
		Key:        key,
		Value:      value,
		Expiration: expiration,
		Duration:   duration,
		Connection: connection,
	}

	source.totalDuration += duration
	source.commands = append(source.commands, &structure)
}

func (source *CacheDataSource) LogCacheMiss(connection, action string, key string, value string, misses int, duration float32, expiration float32) {
	structure := cacheDataStructure{
		Key:        key,
		Value:      value,
		Expiration: expiration,
		Duration:   duration,
		Connection: connection,
	}
	if misses == 1 {
		structure.Type = fmt.Sprintf("MISS %s", action)
	} else {
		structure.Type = fmt.Sprintf("MISSES [%d] %s", misses, action)
	}

	source.cacheReads++
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
