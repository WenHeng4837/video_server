package taskrunner

import (
	"errors"
	"log"
	"os"
	"sync"
	"video_server/scheduler/dbops"
	//"video_server/scheduler/ossops"
)

//流程：api->video_id->mysql
//dispatcher->mysql->video_id->datachannel
//executor->datachannel->video_>delete videos
//真正删除视频
func deleteVideo(vid string) error {
	err := os.Remove(VIDEO_PATH + vid)

	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error: %v", err)
		return err
	}
	/*
		ossfn :="videos/"vid
		bn := "wulihui"
		ok := ossops.DeleteObject(ossfn,bn)
		if ok !=nil{
		log.Printf("Deleting video error,oss operation failed")
		return errors.New("Deleting video error")
		}
	*/
	return nil
}

//
func VideoClearDispatcher(dc dataChan) error {
	//写的时候一定要在外面实际传进来，这里为了简单就直接写60
	//几秒后删除，先写入数据库表，60后删除，延迟删除
	res, err := dbops.ReadVideoDeletionRecord(60)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}
	//如果取出来没有任何结果
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}

	for _, id := range res {
		dc <- id
	}

	return nil
}

//
func VideoClearExecutor(dc dataChan) error {
	//定义一个map
	errMap := &sync.Map{}
	var err error

forloop:
	for {
		select {
		case vid := <-dc:
			//为了并发考虑，这里每一个都起个新的goroutine
			//z这里不直接用vid而使用id是因为你使用go func是一个闭包，当你直接使用vid是会拿到顺时状态，而不会将状态保存，只有你将参数传进来才会保存状态作为一个完整close执行
			go func(id interface{}) {
				//id要转实际类型，先删实际文件
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				//删除表的记录
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break forloop
		}
	}

	errMap.Range(func(k, v interface{}) bool {
		//如果有错直接停止返回
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})

	return err
}
