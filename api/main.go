package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"video_server/api/session"
)

//继承*httprouter.Router，他也是继承了http.Handler
type middleWareHandler struct {
	r *httprouter.Router
}

//生成一个http.Handler，构造函数
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

//middleWareHandler还要继承http.Handler的w http.ResponseWriter, r *http.Request
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

//流程：handler->validation(1.request,2.user)->business lodic->response.
/*
1.data model
2. error handling
*/
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	//创建用户,用户注册
	router.POST("/user", CreateUser)
	//用户登录
	router.POST("/user/:username", Login)
	//用户信息
	router.GET("/user/:username", GetUserInfo)
	//用户资源，增加一个视频(有问题)
	router.POST("/user/:username/videos", AddNewVideo)
	//用户资源,所有视频
	router.GET("/user/:username/videos", ListAllVideos)
	//删除一个用户资源，视频
	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
	//显示评论
	router.POST("/videos/:vid-id/comments", PostComment)
	//提交一个评论
	router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

//启动时把所有用户session加到map
func Prepare() {
	session.LoadSessionsFromDB()
}
func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	//监听
	http.ListenAndServe(":8000", mh)
}
