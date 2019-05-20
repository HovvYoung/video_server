package dbops

import (
	"sync"
	"log"
	"strconv"
	"database/sql"
	"github.com/hovvyoung/video_server/api/defs"
)

func InserSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmIns, err := dbConn.Prepare("INSERT INTO sessions (sessions_id, TTL, login_name) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}

	defer stmIns.Close()
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare("SELECT TTL, login_name FROM sessions WHERE sessions_id = ?")
	if err != nil {
		return nil, err
	}

	var ttl string
	var uname string 
	stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows{
		return nil, err
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err != nil{
		ss.TTL = res  //convert to int64
		ss.Username = uname
	}else{
		return nil, err
	}

	defer stmtOut.Close()
	return ss, nil
}


func RetrieveAllSession() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query()

	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() {
		var id string
		var ttlstr string
		var login_name string

		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
        	log.Printf("retrive sessions error: %s", err)
        	break
        }

        if ttl, err1 := strconv.ParseInt(ttlstr, 10, 64); err1 != nil {
        	ss := &defs.SimpleSession{TTL:ttl, Username:login_name}
        	m.Store(id, ss)
        	log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
        }
	}

	return m, nil
}

func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}

	return nil
}