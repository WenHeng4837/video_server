package main

//import (
//	"github.com/aliyun/aliyun-oss-go-sdk/oss"
//	"log"
//	 "video_server/vendro/wulihui/config"
//)
//阿里云负载均衡配置
//var EP string
//var AK string
//var SK string
//
//func init(){
//	AK=""//你的ak
//	SK=""//你的sk
//	EP=config.GetOssAddr()
//}
////上传
//func UploadToOss(filename string ,path string ,bn string) bool {
//	client ,err := oss.New(EP,AK,SK)
//	if err != nil{
//		log.Printf("Init oss service error :%s",err)
//		return false
//	}
//	bucket ,err := client.Bucket(bn)
//	if err !=nil{
//		log.Printf("Getting bucket error:%s",err)
//		return false
//	}
//	err = bucket.UploadFile(filename,path,500*1024,oss.Routines(200))//并发传
//	if err != nil{
//		log.Printf("Uploading object error :%s",err)
//		return false
//	}
//	return true
//}
