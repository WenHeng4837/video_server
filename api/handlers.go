package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/session"
	"video_server/api/utils"
)

//注册用户
func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//main函数那边router.POST("/user",CreateUser)是post方法，所以这里有个body
	res, _ := ioutil.ReadAll(r.Body)
	//声明实体变量用户名密码
	ubody := &defs.UserCredential{}
	//json.Unmarshal()将res中body反序列化赋值给ubody
	if err := json.Unmarshal(res, ubody); err != nil {
		//这里w是传进来的，错误码和信息是后台定义好的全局变量直接在这里传进去，有点像枚举
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		//sendErrorResponse（）已经将这些信息写入io去自动返回了，所以这里不用再处理返回了
		return
	}
	//注册创建用户
	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	//注册后登陆创建session
	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}
	if resp, err := json.Marshal(su); err != nil {
		//内部错误时
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		//成功时返回sessionId给前端
		sendNormalResponse(w, string(resp), 201)
	}
}

//登录
func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}

	uname := p.ByName("username")
	log.Printf("Login url name:%s", uname)
	log.Printf("Login body name:%s", ubody.Username)
	if uname != ubody.Username {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	log.Printf("%s", ubody.Username)
	//从数据库里拿出来的密码
	pwd, err := dbops.GetUserCredential(ubody.Username)
	log.Printf("Login pwd:%s", pwd)
	//传过来的密码
	log.Printf("Login body pwd:%s", ubody.Pwd)

	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		sendErrorResponse(w, defs.ErrorNotAuthUser)
		return
	}

	//验证通过就给这个用户加sessionID
	id := session.GenerateNewSessionId(ubody.Username)
	si := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

//获取用户信息
func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	////用户校验
	if !ValidateUser(w, r) {
		log.Printf("Unathorized user \n")
		return
	}

	uname := p.ByName("username")
	//视频里时GetUser（）
	u, err := dbops.GetUser(uname)
	if err != nil {
		log.Printf("Error in GetUserInfo ：%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	ui := &defs.UserInfo{Id: u.Id}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		//登录名在登录时就放在cookie里，所以这里只需要返回主键，密码没啥作用
		sendNormalResponse(w, string(resp), 200)
	}
}

//新增一个视频
func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//用户校验
	if !ValidateUser(w, r) {
		log.Printf("Unathorized user \n")
		return
	}
	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal(res, nvbody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	log.Printf("Author id:%d,name:%s\n", nvbody.AuthorId, nvbody.Name)
	if err != nil {
		log.Printf("Error in AddNewVideo:%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	if resp, err := json.Marshal(vi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

//显示所有视频
func ListAllVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//用户校验
	//if !ValidateUser(w,r){
	//	return
	//}
	uname := p.ByName("username")
	vs, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Error in ListAllVideos:%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	vsi := &defs.VideosInfo{Videos: vs}
	if resp, err := json.Marshal(vsi); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

//删除一个视频
func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//用户校验
	if !ValidateUser(w, r) {
		log.Printf("delete video")
		return
	}
	vid := p.ByName("vid-id")
	err := dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("Error in DeleteVideo:%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	//删除数据库顺便删除阿里云
	//go utils.SendDeleteVideoRequest(vid-id)
	sendNormalResponse(w, "", 204)
}

//提交一个评论
func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//用户校验
	if !ValidateUser(w, r) {
		log.Printf("PostComment comment")
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	cbody := &defs.NewComment{}
	if err := json.Unmarshal(reqBody, cbody); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	vid := p.ByName("vid-id")
	if err := dbops.AddNewComments(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("Error in PostComment:%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
	} else {
		sendNormalResponse(w, "ok", 201)
	}
}

//
func ShowComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//用户校验
	if !ValidateUser(w, r) {
		log.Printf("ShowComment comment")
		return
	}
	vid := p.ByName("vid-id")
	cm, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("Eror in ShowComments:%s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	cms := &defs.Comments{Comments: cm}
	if resp, err := json.Marshal(cms); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}
