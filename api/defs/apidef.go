package defs

//实体，有点像java实体类
//requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}
type UserSession struct {
	Username  string `json:"user_name"`
	SessionId string `json:"session_id"`
}
type UserInfo struct {
	Id int `json:"id"`
}

//data model
type User struct {
	Id        int
	LoginName string
	Pwd       string
}
type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}
type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}
type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}
type Comments struct {
	Comments []*Comment `json:"comments"`
}

//response,在用户登录时后创建session的返回
type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

//data model
//视频表的返回实体，不过这里不需要返回创建时间，因为意义不大
type VideoInfo struct {
	Id           string `json:"id"`
	AuthorId     int    `json:"author_id"`
	Name         string `json:"name"`
	DisplayCtime string `json:"display_ctime"`
}

//评论表返回实体
type Comment struct {
	Id      string `json:"id"`
	VideoId string `json:"video_id"`
	//作者名字
	Author  string `json:"author"`
	Content string `json:"content"`
}

//定义一个结构体,session表
type SimpleSession struct {
	Username string //login name
	TTL      int64  // 用来检查是否过期的
}
