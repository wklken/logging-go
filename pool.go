package logging

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

type logEntryPool struct {
	pool sync.Pool
}

func newLogEntryPool() *logEntryPool {
	return &logEntryPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &log.Entry{
					// Logger: logger,
					// Default is three fields, plus one optional.  Give a little extra room.
					Data: make(log.Fields, 6),
				}
			},
		},
	}
}

func (p *logEntryPool) Get(logger *log.Logger) *log.Entry {
	entry := p.pool.Get().(*log.Entry)
	entry.Logger = logger
	return entry
}

func (p *logEntryPool) Put(e *log.Entry) {
	// TODO: clean, should make?
	// reference: https://github.com/sirupsen/logrus/pull/796/files
	// e.Data = make(log.Fields, 6)
	e.Data = map[string]interface{}{}

	p.pool.Put(e)
}
