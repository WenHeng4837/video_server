package main

import (
	"encoding/json"
	"io"
	"net/http"
	"video_server/api/defs"
)

//错误返回
func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	//defs.ErrorResponse里面的有点想状态码的东西
	w.WriteHeader(errResp.HttpSC)
	//通过json.Marshal（）将错误信息返回
	resStr, _ := json.Marshal(&errResp.Error)
	//写进去io里返回
	io.WriteString(w, string(resStr))
}

//这里校验格式可能HttpSC以及错误信息会因前端不同情况而变化，所以这里作为参数传进来
//这里不一定是错误返回
func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	io.WriteString(w, resp)
}
