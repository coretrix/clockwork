package datasource

type CommandLoggerDataSourceInterface interface {
	DataSource
	CommandLoggerInterface
}

type CommandLoggerInterface interface {
	LogCommand(command string, key string, duration float32)
}

type redisStructure struct {
	Command    string            `json:"command"`
	Parameters map[string]string `json:"parameters"`
	Duration   float32           `json:"duration"`
}

type RedisDataSource struct {
	commands []interface{}
}

func (source *RedisDataSource) LogCommand(command string, key string, duration float32) {
	structure := redisStructure{
		Command:    command,
		Parameters: map[string]string{"key": key},
		Duration:   duration,
	}

	source.commands = append(source.commands, &structure)
}

func (source *RedisDataSource) Resolve(dataBuffer *DataBuffer) {
	dataBuffer.RedisCommands = source.commands
}
