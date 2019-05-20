package main

import (
	"log"
)

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
	cl.bucket <- 1
	return true  //成功拿到token返回true
}

func (cl * ConnLimiter) ReleaseConn() {
	c := <- cl.bucket
	log.Printf("New connection coming: %d", c)
}