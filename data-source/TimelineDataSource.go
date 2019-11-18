package dataSource

import (
	"time"
)

type TimelineLoggerDataSourceInterface interface {
	DataSource
	TimelineLoggerInterface
}

type TimelineLoggerInterface interface {
	StartEvent(event string, description string)
	EndEvent(event string)
}

type timelineStructure = struct {
	Start float64 `json:"start"`
	End float64 `json:"end"`
	Duration float64 `json:"duration"`
	Description string `json:"description"`
}

type TimelineDataSource struct {
	commands map[string]interface{}
	startTime map[string]time.Time
	description map[string]string
	data []string
}

func (source *TimelineDataSource) StartEvent(event string, description string)  {
	if len(source.startTime) == 0 {
		source.startTime = make(map[string]time.Time)
	}

	if len(source.commands) == 0 {
		source.commands = make(map[string]interface{})
	}

	if len(source.description) == 0 {
		source.description = make(map[string]string)
	}

	source.description[event] = description
	source.startTime[event] = time.Now()
}

func (source *TimelineDataSource) EndEvent(event string)  {
	start := MicroTime(source.startTime[event])
	end := MicroTime(time.Now())

	source.commands[event] = timelineStructure {
		Start: start,
		End: end,
		Duration: (end * 1000) - (start * 1000),
		Description: source.description[event],
	}
}

func (source *TimelineDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.TimelineData = source.commands
}
