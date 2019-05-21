package taskrunner

import (
	"time"
	//"log"

) 

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker {
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <- w.ticker.C:	//时间到了系统会往ticker.C这个channel发送东西
			go w.runner.StartAll()
		}
	}
}

func Start() {
	// Start video file cleaning. We read 3 rows per time
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)
	go w.startWorker()
}