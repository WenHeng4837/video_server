package main

import (
	"net/http"
	"video_server/api/defs"
	"video_server/api/session"
)

//流程：先检查用户是否存在再检查session是否合法

//http自定义header,X开头
var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

//检查用户session是否合法session ,不是的话返回false，是就是true
func validateUserSession(r *http.Request) bool {
	//sid直接false
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	//是否过期
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	//不过期用户名加进来
	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

//用户校验
func ValidateUser(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_UNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return false
	}
	return true
}
