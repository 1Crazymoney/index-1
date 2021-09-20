package run

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/jlog"
	"github.com/memocash/server/db/server"
)

type Queue struct {
	Port   uint
	Server *server.Server
	Error  error
}

func (q *Queue) Start() error {
	q.Server = server.NewServer(q.Port)
	jlog.Logf("Starting queue server on port: %d\n", q.Port)
	go func() {
		err := q.Server.Run()
		q.Error = jerr.Get("error queue server ended", err)
	}()
	return nil
}

func (q *Queue) End() {
	q.Server.Stop()
}

func NewQueue(port uint) *Queue {
	return &Queue{
		Port: port,
	}
}
