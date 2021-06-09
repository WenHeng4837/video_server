package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

//测试文件
var tempvid string //全局变量

func clearTables() {
	//truncate删除此表所有数据，不产生二进制日志，无法恢复数据，速度极快
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	//先初始化，清空表的数据
	clearTables()
	//测试
	m.Run()
	//再初始化。清空表的数据
	clearTables()
}

//在这里依次测试
func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

//测试增加用户
func testAddUser(t *testing.T) {
	err := AddUserCredential("wulihui", "123")
	if err != nil {
		t.Errorf("Error of AddUser:%v", err)
	}
}

//测试得到用户
func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("wulihui")
	if pwd != "123" || err != nil {
		t.Errorf("Error of GetUser")
	}
}

//测试删除用户
func testDeleteUser(t *testing.T) {
	err := DeleteUser("wulihui", "123")
	if err != nil {
		t.Errorf("Error of DeleteUser:%v", err)
	}
}

//这个函数用于增加获取删除用户后看还能不能获取到该用户有没有真正删除掉
func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("wulihui")
	if err != nil {
		t.Errorf("Error of RegetUser:%v", err)
	}
	if pwd != "" {
		t.Errorf("Deleting user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	//这里先清除数据，再加个用户才能测视频
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

//测试增加视频
func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

//测试获取视频
func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

//测试删除视频
func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

//测试用于增加获取删除视频后看还能不能获取到该视频有没有真正删除掉
func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil {
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

//测试评论
func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddCommnets", testAddComments)
	t.Run("ListComments", testListComments)
}

//测试增加评论
func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	err := AddNewComments(vid, aid, content)

	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

//测试同个视频所有获取评论
func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	//转换当前时间
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}
	//遍历
	for i, ele := range res {
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}
