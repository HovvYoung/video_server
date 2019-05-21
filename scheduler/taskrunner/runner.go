package taskrunner

import (
)

type Runner struct {
	Controller controlChan
	Error controlChan	//正常信息和错误信息分开， 便于维护
	Data dataChan
	dataSize int
	longLived bool  //是否回收
	Dispatcher fn 
	Executor fn
}

func NewRunner(size int, longlived bool, d fn, e fn) *Runner {
	return &Runner {
		Controller: make(chan string, 1),	//设置成非阻塞
		Error: make(chan string, 1),
		Data: make(chan interface{}, size),
		longLived: longlived,
		dataSize: size,
		Dispatcher: d,
		Executor: e,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	//生产者消费者模型， 注意这里是阻塞的
	for {
		select {
		case c:=<- r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			}

			if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE 
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e :=<- r.Error:
			if e == CLOSE {
				return
			}
		default:
		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}