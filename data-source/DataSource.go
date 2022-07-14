package dataSource

import "time"

type DataSource interface {
	Resolve(dataBuffer *DataBuffer)
}

func MicroTime(timer time.Time) float64 {
	return float64(timer.Unix()) + float64(timer.Nanosecond())/float64(time.Second)
}
