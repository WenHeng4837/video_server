package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"video_server/api/defs"
	"video_server/api/utils"
)

//添加用户凭证也就是添加用户
func AddUserCredential(loginName string, pwd string) error {
	//预编译
	stmtIns, err := dbConn.Prepare("INSERT INTO users(login_name,pwd) VALUES(?,?)")
	if err != nil {
		return err
	}
	//执行
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	//关闭
	//defer语句的作用是不管程序是否出现异常，均在函数退出时自动执行相关代码
	defer stmtIns.Close()
	return nil
}

//得到用户凭证,也就是得到密码
func GetUserCredential(loginName string) (string, error) {
	stmtOUt, err := dbConn.Prepare("select pwd from users where login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}
	var pwd string
	err = stmtOUt.QueryRow(loginName).Scan(&pwd)
	//如果不出错但是查询不到数据也会返回一个Object叫Row,这个Scan会把错误带出来这个错误叫NoRows
	//ErrNoRows其实不是一个真正的错误结果，只是没有结果按照一个错误来返回了
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOUt.Close()
	return pwd, nil
}

//得到用户主键
func GetUser(loginName string) (*defs.User, error) {
	stmtOUt, err := dbConn.Prepare("select id,login_name,pwd from users where login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	user := &defs.User{}
	err = stmtOUt.QueryRow(loginName).Scan(&user.Id, &user.LoginName, &user.Pwd)
	//如果不出错但是查询不到数据也会返回一个Object叫Row,这个Scan会把错误带出来这个错误叫NoRows
	//ErrNoRows其实不是一个真正的错误结果，只是没有结果按照一个错误来返回了
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer stmtOUt.Close()
	return user, nil
}

//删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("delete from users where login_Name= ? and pwd =?")
	if err != nil {
		log.Printf("DeleteUser error : %s", err)
		stmtDel.Exec(loginName, pwd)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//添加视频
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	//创建一个uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	//按照格式转成对应时间字符串
	ctime := t.Format("Jan 02 2006,15:04:05") //M,D,Y,HH:MM:SS
	stmIns, err := dbConn.Prepare("insert into video_info(id,author_id,name,display_ctime) values(?,?,?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	defer stmIns.Close()
	return res, nil
}

//查询视频
func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")
	var aid int
	var dct string
	var name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}
	return res, nil
}

//删除视频
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//评论不能删
//增加评论
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

//查评论
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(` SELECT comments.id, users.Login_name, comments.content FROM comments
		INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)
		ORDER BY comments.time DESC`)

	var res []*defs.Comment

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return res, err
		}
		//这name是从用户表里来的
		c := &defs.Comment{Id: id, VideoId: vid, Author: name, Content: content}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil
}

//查所有视频
func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`select video_info.id ,video_info.author_id,video_info.name,video_info.display_ctime from video_info
 		where users.id in (select users.id from users where users.login_name=?)=video_info.author_id and 
 		video_info.display_ctime > FROM_UNIXTIME(?) AND video_info.display_ctime <= FROM_UNIXTIME(?)
 		ORDER BY video_info.display_ctime DESC`)
	var res []*defs.VideoInfo
	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var id string
		var authorid int
		var name string //视频名字
		var displayctime string
		if err := rows.Scan(&id, &authorid, &name, &displayctime); err != nil {
			return res, err
		}
		c := &defs.VideoInfo{Id: id, AuthorId: authorid, Name: name, DisplayCtime: displayctime}
		res = append(res, c)
	}
	defer stmtOut.Close()
	return res, nil

}
