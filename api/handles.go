package main

import (
	"io"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/hovvyoung/video_server/api/defs"
	"github.com/hovvyoung/video_server/api/dbops"
	"github.com/hovvyoung/video_server/api/session"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	 //post
	 res, _ := ioutil.ReadAll(r.Body)
	 ubody := &defs.UserCredential{}

	 if err := json.Unmarshal(res, ubody); err != nil {
	 	sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
	 	return
	 }

	 if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil{
	 	sendErrorResponse(w, defs.ErrorDBError)
	 }

	 id := GenerateNewSessionId(ubody.Username)
	 su := &defs.SignedUp{Success:true, SessionId: id}

	 if resp, err := json.Marshal(su); err != nil {
	 	sendErrorResponse(w, defs.ErrorInternalFaults)
	 }else {
	 	sendNormalResponse(w, resp, 201)
	 }
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}