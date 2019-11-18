package dataSource

import (
	"runtime"
	"time"
)

const (
	debugLevel = "debug"
)

type LoggerDataSourceInterface interface {
	DataSource
	LoggerInterface
}

type LoggerInterface interface {
	LogDebugSlice(name string, debugData []string)
	LogDebugMap(name string, debugData map[string]string)
	LogDebugString(name string, debugData string)
}

type loggerStructure = struct {
	Message string `json:"message"`
	Context interface{} `json:"context"`
	Time float64 `json:"time"`
	Level string `json:"level"`
	File string `json:"file"`
	Line int `json:"line"`
}

type LoggerDataSource struct {
	commands []interface{}
}

func (source *LoggerDataSource) LogDebugSlice(name string, debugData []string)  {
	elementMap := make(map[string]string)
	for i := 0; i < len(debugData); i +=2 {
		elementMap[debugData[i]] = debugData[i+1]
	}
	elementMap["__type__"] = "array"

	_, file, line, _ := runtime.Caller(1)

	structure := loggerStructure{
		Message: name,
		Context: debugData,
		Time: MicroTime(time.Now()),
		Level: debugLevel,
		File: file,
		Line: line,
	}
	source.commands = append(source.commands, &structure)
}

func (source *LoggerDataSource) LogDebugMap(name string, debugData map[string]string)  {
	debugData["__type__"] = "array"
	_, file, line, _ := runtime.Caller(1)

	command := loggerStructure{
		Message: name,
		Context: debugData,
		Time: MicroTime(time.Now()),
		Level: debugLevel,
		File: file,
		Line: line,
	}
	source.commands = append(source.commands, &command)
}

func (source *LoggerDataSource) LogDebugString(name string, debugData string)  {
	elementMap := make(map[string]string)
	elementMap["__type__"] = "array"
	elementMap["__string__"] = debugData

	_, file, line, _ := runtime.Caller(1)

	command := loggerStructure{
		Message: name,
		Context: elementMap,
		Time: MicroTime(time.Now()),
		Level: debugLevel,
		File: file,
		Line: line,
	}
	source.commands = append(source.commands, &command)
}

func (source *LoggerDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.Log = source.commands
}
