package main

//流控防止恶意不断攻击访问导致带宽被消耗完处于无法访问状态
import (
	"log"
)

//连接流控结构体
type ConnLimiter struct {
	//连接数
	concurrentConn int
	//bucket数量,一个channel维护着一个长链接用户,chan关键字
	bucket chan int
}

//ConnLimiter构造函数
func NewConnLimiter(cc int) *ConnLimiter {
	//引用返回
	return &ConnLimiter{
		concurrentConn: cc,
		//创造一个新的chan，跟连接一样数量才能同步
		//当新的request进来就写入chan，走了再释放
		bucket: make(chan int, cc),
	}
}

//用来判断token是否拿到
func (cl *ConnLimiter) GetConn() bool {
	//满了就返回错误
	if len(cl.bucket) >= cl.concurrentConn {
		log.Printf("Reached the rate limitation.")
		return false
	}
	//没满就写进一个1返回true表示拿到token
	cl.bucket <- 1
	return true
}

//释放token
func (cl *ConnLimiter) ReleaseConn() {
	//把写进去的token拿出来
	c := <-cl.bucket
	log.Printf("New connction coming: %d", c)
}
