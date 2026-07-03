package engine

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	engine *Engine
	ticker *time.Ticker
	done   chan bool
}

func NewScheduler(engine *Engine, interval time.Duration) *Scheduler {
	return &Scheduler{
		engine: engine,
		ticker: time.NewTicker(interval),
		done:   make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("Starting Recurring Transaction Scheduler...")
	go func() {
		for {
			select {
			case <-s.done:
				return
			case t := <-s.ticker.C:
				log.Printf("Scheduler tick at %v", t)
				if err := s.engine.ProcessDueTransactions(ctx, t); err != nil {
					log.Printf("Error processing due transactions: %v", err)
				}
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	s.ticker.Stop()
	s.done <- true
}
