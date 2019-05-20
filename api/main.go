package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type middleWareHandle struct {
	r *httprouter.Router
}

func NewMiddleWareHandle(r *httprouter.Router) http.Handler {
	m := middleWareHandle{}
	m.r = r
	return m
}

func (m middleWareHandle) ServeHttp(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)
	
	m.r.ServeHttp(w, r)
}

func RegisterHandles() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)

	router.POST("/user/:user_name", Login)

	return router
}

func main() {
	r := RegisterHandles()
	mh := NewMiddleWareHandle(r)
	http.ListenAndServe(":8000", mh)
}