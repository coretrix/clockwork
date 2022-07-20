package datasource

import (
	"fmt"
	"sync"
)

const typeTable = "table"
const typeCounters = "counters"

type UserDataSourceInterface interface {
	DataSource
	UserDataInterface
}

type UserDataInterface interface {
	SetTitle(title string)
	SetShowAs(showAs string)
	LogTable(data map[string]interface{}, title string, labels map[string]string)
}

type UserDataDataSource struct {
	mutex       *sync.Mutex
	rowsByTitle map[string]map[string]interface{}
	counter     int
	Title       string                   `json:"title"`
	Data        []map[string]interface{} `json:"data"`
	ShowAs      string                   `json:"showAs"` // Describes how the data should be presented ("counters" or "table")
}

//LogTable
//labels:  Map of human-readable labels for the data contents
//showAs:  Describes how the data should be presented ("counters" or "table")
func (source *UserDataDataSource) LogTable(data map[string]interface{}, title string, labels map[string]string) {
	source.lock()
	defer source.unlock()

	if source.rowsByTitle == nil {
		source.rowsByTitle = map[string]map[string]interface{}{}
	}

	if source.rowsByTitle[title+"_table"] == nil {
		source.rowsByTitle[title+"_counters"] = map[string]interface{}{}
		source.rowsByTitle[title+"_counters"]["__meta"] = map[string]interface{}{
			"showAs": typeCounters,
			"labels": map[string]string{"Total": "Total"},
			"title":  title,
		}

		source.rowsByTitle[title+"_table"] = map[string]interface{}{}
		source.rowsByTitle[title+"_table"]["__meta"] = map[string]interface{}{
			"showAs": typeTable,
			"labels": labels,
			"title":  title,
		}
	}

	source.counter++
	source.rowsByTitle[title+"_counters"]["Total"] = source.counter
	source.rowsByTitle[title+"_table"][fmt.Sprint(source.counter)] = data
}

func (source *UserDataDataSource) SetShowAs(showAs string) {
	source.lock()
	defer source.unlock()

	source.ShowAs = showAs
}

func (source *UserDataDataSource) SetTitle(title string) {
	source.lock()
	defer source.unlock()

	source.Title = title
}

func (source *UserDataDataSource) Resolve(dataBuffer *DataBuffer) {
	source.lock()
	defer source.unlock()

	if dataBuffer.UserData == nil {
		dataBuffer.UserData = make([]map[string]interface{}, 0)
	}

	data := map[string]interface{}{
		"__meta": map[string]string{
			"title":  source.Title,
			"showAs": source.ShowAs,
		},
	}

	for k, v := range source.rowsByTitle {
		data[k] = v
	}
	dataBuffer.UserData = append(dataBuffer.UserData, data)
}

func (source *UserDataDataSource) lock() {
	if source.mutex == nil {
		source.mutex = &sync.Mutex{}
	}
	source.mutex.Lock()
}
func (source *UserDataDataSource) unlock() {
	if source.mutex == nil {
		panic("lock first")
	}
	source.mutex.Unlock()
}
