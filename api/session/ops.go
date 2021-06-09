package session

import (
	"sync"
	"time"
	_ "time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

//不用redis是因为当你增加一个模块或者一个东西的时候势必系统复杂度会增加，但是增加的复杂度大于它业务上带来的好处就没有必要
//综合考虑map，这个项目用户数据量有限，map在每个节点都会缓存她所有的数据
//并发读写上表现可以
//全局map
var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

//获取当前时间并且用毫秒数
func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

//删除过期session
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

//从db拉取session到缓存
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	//内置函数，遍历时把id当k，然后定义interface{}类型的v，在函数里实际运用时才指定真正类型为defs包里的SimpleSession
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

//新注册产生sessionID
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	//过期时间是当前登录时间加上30分钟
	ct := nowInMilli()
	ttl := ct + 30*60*1000 // Severside session valid time: 30 min
	ss := &defs.SimpleSession{Username: un, TTL: ttl}
	//注册产生session后面可能要登录，所以先把session存进全局map缓存里再插入数据库
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, un)

	return id
}

//校验的时候session过期或者不过期返回登录状态
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid)
	ct := nowInMilli()
	if ok {
		//ct1 := nowInMilli()
		//全局map的值是interface{}型所以这里要转换成实际类型
		if ss.(*defs.SimpleSession).TTL < ct {
			//过期
			deleteExpiredSession(sid)
			return "", true
		} else {
			ss, err := dbops.RetrieveSession(sid)
			//一个是真的有错误一个是没有session
			if err != nil || ss == nil {
				return "", true
			}
			if ss.TTL < ct {
				deleteExpiredSession(sid)
				return "", true
			}
		}
		//不过期
		sessionMap.Store(sid, ss)
		return ss.(*defs.SimpleSession).Username, false
	}
	//如果执行出错就返回这个
	return "", true
}
