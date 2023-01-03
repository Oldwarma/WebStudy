package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

//超时的中间件
func TimeoutHandler(fun ControllerHandler, d time.Duration) ControllerHandler {
	//使用函数回调
	return func(c *Context) error {
		// 这个 channal 负责通知结束
		finish := make(chan struct{}, 1)
		// 这个 channel 负责通知 panic 异常
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
		defer cancel()

		c.request.WithContext(durationCtx)
		go func() {
			// 这里增加异常处理
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 这里做具体的业务
			fun(c)

			finish <- struct{}{}
		}()
		//执行业务逻辑后操作
		select {
		case p := <-panicChan:
			log.Println(p)
			c.responseWriter.WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.responseWriter.Write([]byte("time out"))
		}
		return nil
	}
}
