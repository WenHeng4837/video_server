package main

//这边大前端后台主要做两件事：
//1.请求转发（防止跨域问题）
//2.把返回的东西渲染到前端展示
//这里还要用这个的原因是因为虽然这里渲染的还是页面但是一样是前台的request到后台，
//后台接受处理再返回给前台，只不过在web里面返回的都是整个页面而不会是json的一个消息或者是视频的一个文件流
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()
	//页面首次跳转登陆进去
	router.GET("/", homeHandler)
	//不是首次，加载提交表单跳转到加载的初始登陆页面
	router.POST("/", homeHandler)
	//用户登录后个人页面
	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	//api转发
	router.POST("/api", apiHandler)
	//proxy代理
	//aliyun
	//router.Get("/videos/:vid-id",proxyVideoHandler)
	router.POST("/upload/:vid-id", proxyUploadHandler)
	//静态文件
	router.ServeFiles("/statics/*filepath", http.Dir("./template/"))

	return router
}
func main() {
	r := RegisterHandler()
	http.ListenAndServe(":8080", r)
}
