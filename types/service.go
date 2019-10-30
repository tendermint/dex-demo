package types

import (
	tmlog "github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/dex-demo/pkg/log"
)

type LocalConsumer struct {
	queue    Backend
	quitCh   chan bool
	handlers []EventHandler
	lgr      tmlog.Logger
}

func NewLocalConsumer(queue Backend, hdlrs []EventHandler) *LocalConsumer {
	return &LocalConsumer{
		queue:    queue,
		handlers: hdlrs,
		lgr:      log.WithModule("local-consumer"),
	}
}

func (s *LocalConsumer) Start() {
	go func() {
		for {
			select {
			case <-s.quitCh:
				return
			default:
				item := s.queue.Consume()
				s.handleItem(item)
			}
		}
	}()
}

func (s *LocalConsumer) Stop() {
	s.quitCh <- true
}

func (s *LocalConsumer) handleItem(item interface{}) {
	for _, hdlr := range s.handlers {
		if err := hdlr.OnEvent(item); err != nil {
			s.lgr.Error("error consuming queue item", "err", err.Error())
		}
	}
}
