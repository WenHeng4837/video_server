package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	_ "video_server/vendro/wulihui/config"
)

//最终切换到首页的那个模板对象
type HomePage struct {
	Name string
}

//登录状态的模板对象
type UserPage struct {
	Name string
}

//未登录
func homeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//获取用户名和session
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	//游客，去注册
	if err1 != nil || err2 != nil {
		p := &HomePage{Name: "wulihui"}
		t, e := template.ParseFiles("./template/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}
		//把模板和变量一起渲染进去
		t.Execute(w, p)
		return
	}
	//用户，不是游客并且点了登录按钮提交用户密码直接去登录页面
	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		//重定向到userhome
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}

}

//登录
func userHomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	//没有登录或者第一次到这个页面
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	//用户通过homepage后从表单进来的，这里前端调用别的方法去判断是否合法，所以这里不用判断
	fname := r.FormValue("username")
	var p *UserPage
	if len(cname.Value) != 0 {
		//已经登录从cookie里拿出来
		p = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		//如果没有登录尝试从表单提交的记录去读
		p = &UserPage{Name: fname}
	}
	t, e := template.ParseFiles("./template/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}
	t.Execute(w, p)
}

//api请求转发
func apiHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//先校验是否是post方式
	if r.Method != http.MethodPost {
		//错误请求，json.Marshal：将数据编码为json格式
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(re))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	//赋值过去
	if err := json.Unmarshal(res, apibody); err != nil {
		//body不符合要求
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(re))
		return
	}
	request(apibody, w, r)
	defer r.Body.Close()
}

//跨域请求代理(文件)
/*
//阿里云
func proxyVideoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//代理9000这个端口，这个端口跟视频有关的
	u, _ := url.Parse("http://"+config.GetLbAddr()+":9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
*/
func proxyUploadHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//代理9000这个端口，这个端口跟视频有关的
	u, _ := url.Parse("http://127.0.0.1:9000/")
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
