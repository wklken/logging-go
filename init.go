package logging

var LogEntryPool *logEntryPool

func init() {
	LogEntryPool = newLogEntryPool()
}
