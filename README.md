# video_server

### stream模块
rateLimiter 限流

```
// streamserver/main.go
type middleWareHandler struct {
	r *httprouter.Router    // 路由器
	l *ConnLimiter  // 限流器
}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//if not get token
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests,"Too many requests.")  //SCode: 429
		return
	}
    // 拿到token才继续
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 2)
	http.ListenAndServe(":9000", mh)
}
```
```
// streamserver/limiter.go
type ConnLimiter struct {
	concurrentConn int  //max connections
	bucket chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter {
		concurrentConn: cc,
		bucket: make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	//当前bucket里运行中的个数 > max
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reach the rate limitation.")
		return false
	}
	cl.bucket <- 1 // 一个连接从bucket拿一个token
	return true  //成功拿到token返回true
}

func (cl * ConnLimiter) ReleaseConn() {
	c := <- cl.bucket // 连接结束还回去
	log.Printf("New connection coming: %d", c)
}

```

# Scheduler模块
```
Worker
    ticker
  - runner
        Controller: make(chan string, 1),	//设置成非阻塞
        Error: make(chan string, 1),
        Data: make(chan interface{}, size),
        longLived: longlived,
        dataSize: size,
        Dispatcher: d, // d, e两个自定义函数，参数传入
        Executor: e,        
```

```
main.go
func main() {
	//启动taskrunner
	go taskrunner.Start()
	r := RegisterHandlers()

	http.ListenAndServe(":9001", r)
}
```
```
// trmain.go

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
            //时间到了系统会往ticker.C这个channel发送东西,这里就可以取到
            case <- w.ticker.C:	
                go w.runner.StartAll()
        }
    }
}

func Start() {
	// Start video file cleaning. We read 3 rows per time
	r := NewRunner(3, true, VideoClearDispatcher, VideoClearExecutor)
	w := NewWorker(3, r)

    // 可以配置多个不同的定时任务
    r := NewRunner(...)
    w := NewWorker(...)

	go w.startWorker()
}
```

```
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
    // 开始的时候往controChan发送Dispatch信息
    r.Controller <- READY_TO_DISPATCH
    r.startDispatch()
}
```
