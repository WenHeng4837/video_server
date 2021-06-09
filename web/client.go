package main

//前端代理用于转发request给api去处理的
import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	_ "video_server/vendro/wulihui/config"
)

//全局client
var httpClient *http.Client

//先初始化
func init() {
	httpClient = &http.Client{}
}

//请求转发
func request(b *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	/*
			u,_ := url.Parse(b.Url)
			u.Host = config.GetLbAddr() + ":" + u.Port()
			newUrl := u.String()
		//然后把下面所有的b.Url换成newUrl
	*/
	//method方式
	//go语言里很完美不用自动去break
	switch b.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", b.Url, nil)
		req.Header = r.Header
		//httpClient.Do（）这个方法去转发
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, resp)
	case http.MethodPost:
		req, _ := http.NewRequest("POST", b.Url, bytes.NewBuffer([]byte(b.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, resp)
	case http.MethodDelete:
		req, _ := http.NewRequest("Delete", b.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		normalResponse(w, resp)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}

}

//正常回应
func normalResponse(w http.ResponseWriter, r *http.Response) {
	//有错也要算正常反应返回去
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(re))
		return
	}
	//没错
	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}
