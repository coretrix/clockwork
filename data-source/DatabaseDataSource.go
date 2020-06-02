package dataSource

import (
	"fmt"
)

type QueryLoggerDataSourceInterface interface {
	DataSource
	QueryLoggerInterface
}

type QueryLoggerInterface interface {
	LogQuery(model, query string, duration float32, bind []interface{})
}

type mySqlStructure = struct {
	Model      string   `json:"model"`
	Query      string   `json:"query"`
	Duration   float32  `json:"duration"`
	Connection string   `json:"connection"`
	Tags       []string `json:"tags"`
}

type DatabaseDataSource struct {
	commands      []interface{}
	totalDuration float32
}

func (source *DatabaseDataSource) LogQuery(model, query string, duration float32, bind []interface{}) {
	var tags []string

	if duration > 2 {
		tags = append(tags, "slow")
	}

	structure := mySqlStructure{
		Model:      model,
		Query:      query + " [" + fmt.Sprintf("%v", bind) + "]",
		Duration:   duration,
		Connection: "test-connection",
		Tags:       tags,
	}

	source.totalDuration += duration
	source.commands = append(source.commands, &structure)
}

func (source *DatabaseDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.DatabaseQueries = source.commands
	dataBuffer.DatabaseDuration = source.totalDuration
	dataBuffer.DatabaseQueriesCount = len(dataBuffer.DatabaseQueries)
}
