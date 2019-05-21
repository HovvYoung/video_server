package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/hovvyoung/video_server/api/session"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)

	m.r.ServeHTTP(w, r)
}

func RegisterHandles() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)

	router.POST("/user/:username", Login)

	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)

	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

func Prepare() {
	session.LoadSessionsFromDB()	//sessions预处理, Load to local memory
}

func main() {
	r := RegisterHandles()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)
}

/*upload
logout*(not)
ListVideos
ShowComments
PostComment
DeleteVideo*/