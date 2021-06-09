package taskrunner

const (
	//这三个都是controller的消息
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE  = "e"
	//以上两个消息一旦谁发生错误就会触发下面这个
	CLOSE = "c"
	//
	VIDEO_PATH = "./videos/"
)

//Controller的chan
type controlChan chan string

//data的chan,泛型
type dataChan chan interface{}

//一个函数类型的结构fn,这个函数参数是dataChan变量，返回错误:dispacter和Executor
type fn func(dc dataChan) error
