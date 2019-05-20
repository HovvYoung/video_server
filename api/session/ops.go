package session

import (
	"time"	//isExpired?
	"sysn"	//save session
	"github.com/hovvyoung/video_server/api/defs"
)



var sessionMap *sync.Map  //研究一下

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano()/1000000
}

func LoadSessionFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}

	r.Range(func(k, v interface{}) bool{
		ss := v.(*defs.SimpleSession)  //通过断言实现类型转换
		sessionMap.Store(k, ss)
		return true
	})
}

func GenerateNewSessionId(un string) string {
	id,_ := utils.NewUUID()
	ct := time.Now().UnixNano()/1000000	//ms
	ttl := ct + 30*60*1000	//severside session valid time : 30 mins

	ss := &defs.SimpleSession{Username:un, TTL:ttl}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, un)

	return id
}

//if not expired, return username and bool.
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()

	if ok {
		//通过断言实现类型转换
		if ss.(*defs.SimpleSession).TTL < ct {
			deleteExpiredSession(sid)	//time.now has already been beyond TTL.
			return "", true
		}

		return ss.(*defs.SimpleSession).Username, false
	} else {
		//if not in memory, try to retrieve from DB.
		ss, err := dbops.RetrieveSession(sid)
		if err != nil || ss != nil {
			return "", true
		}

		if ss.TTL < ct {
			deleteExpiredSession(sid)
			return "", true
		}

		sessionMap.Store(sid, ss)
		return ss.Username, false
	}

	return "", true
}