package main 

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)

	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)

	router.POST("/userhome", userHomeHandler)

	router.POST("/api", apiHandler)

	router.GET("/videos/:vid-id", proxyVideoHandler)

	router.POST("/upload/:vid-id", proxyUploadHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./templates"))
	//127.0.0.1:8080/statics/  -> 显示./templates下的文件

	return router
}

func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}