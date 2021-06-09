package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//加载前端页面
func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")

	t.Execute(w, nil)
}

//将service 端文件传到客户端
func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid
	//打开视频
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	//将它设置为video的MP4格式，浏览器就按照这个格式解析
	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)
	defer video.Close()
	/*//云存储
	log.Println("Entered the streamHandler")
	targetUrl := "http://wulihui.oss-cn-beijing.aliyuncs.com/videos/"+p.ByName("vid-id")
	http.Redirect(w,r,targetUrl,301)
	*/
}

//将客户端文件视频传到service 端
func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//指定文件上传最大字节数
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	//底层通过调用multipartReader.ReadForm来解析
	//如果文件大小超过maxMemory,则使用临时文件来存储multipart/form中文件数据
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}
	//接受前端名字为file的文件，这里的下划线其实是返回一个handler用来校验的，但是这里校验一般在前端校验的
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	//再把文件里的数据读出来
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
	}
	//名字是放在参数里的
	fn := p.ByName("vid-id")
	//把数据写到文件夹的对应的名字视频里
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	/*//有问题filename
	//上传oss
	ossfn := "videos/"+filename
	path := "./videos/" + filename
	bn := "wulihui"
	ret := UploadToOss(ossfn,path,bn)
	if !ret{
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	//上传成功后从临时文件里移除
	os.Remove(path)
	*/
	//给它返回个正确的response
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
	/*

	 */
}
