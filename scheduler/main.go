package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/hovvyoung/video_server/scheduler/taskrunner"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video-delete-record/:vid-id", vidDelRecHandler)

	return router
}

func main() {
	//c: make(chan int)
	go taskrunner.Start()
	r := RegisterHandlers()
	//<- c  a way to prevent from exit;

	http.ListenAndServe(":9001", r)
}