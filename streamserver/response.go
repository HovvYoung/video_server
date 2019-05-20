package main

import (
	//"encoding/json"
	"io"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)  //WriterHeader 要放在w.Header().Set()之后
	io.WriteString(w, errMsg)
}