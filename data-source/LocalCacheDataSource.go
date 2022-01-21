package dataSource

type UserDataSourceInterface interface {
	DataSource
	UserDataInterface
}

type UserDataInterface interface {
	SetTitle(title string)
	SetShowAs(showAs string)
	Log(data map[string]interface{}, title, showAs string, labels map[string]string)
}

type UserDataDataSource struct {
	Title  string                   `json:"title"`
	Data   []map[string]interface{} `json:"data"`
	ShowAs string                   `json:"showAs"` // Describes how the data should be presented ("counters" or "table")
}

//labels:  Map of human-readable labels for the data contents
//showAs:  Describes how the data should be presented ("counters" or "table")
func (source *UserDataDataSource) Log(data map[string]interface{}, title, showAs string, labels map[string]string) {
	source.Data = append(source.Data, map[string]interface{}{
		"data": data,
		"__meta": map[string]interface{}{
			"showAs": showAs,
			"labels": labels,
			"title":  title,
		},
	})
}

func (source *UserDataDataSource) SetShowAs(showAs string) {
	source.ShowAs = showAs
}

func (source *UserDataDataSource) SetTitle(title string) {
	source.Title = title
}

func (source *UserDataDataSource) Resolve(dataBuffer *DataBuffer) {
	if dataBuffer.UserData == nil {
		dataBuffer.UserData = make([]map[string]interface{}, 0)
	}

	dataBuffer.UserData = append(dataBuffer.UserData, map[string]interface{}{
		"Title":  source.Title,
		"Data":   source.Data,
		"ShowAs": source.ShowAs,
	})
}
