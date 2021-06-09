package taskrunner

//runner测试
import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	//定义一个dispacter
	d := func(dc dataChan) error {
		//往里面写东西
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("Dispatcher sent: %v", i)
		}

		return nil
	}
	//定义一个Executor
	e := func(dc dataChan) error {
		//把chan里面的消息打印出来，forloop保证打印完就跳出不会来回执行
	forloop:
		for {
			select {
			case d := <-dc:
				log.Printf("Executor received: %v", d)
			default:
				break forloop
			}
		}

		return errors.New("Executor")

	}

	runner := NewRunner(30, false, d, e)
	//这里没有直接调用而是起了个go,是因为StartAll()里的startDispatch()是个死循环，如果不在后台给它挂起就会一直break死循环不会走下面的time.Sleep
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
