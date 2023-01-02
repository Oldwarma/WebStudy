package test

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const short = 1 * time.Millisecond

func TestOnce(t *testing.T) {
	//创建截止时间
	d := time.Now().Add(short)
	//创建有截止时间的Context
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	//使用select 监听1s和有截止时间的Context哪个先结束
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("超时")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}

}

func TextError(t *testing.T) {

	//durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	//// 这里记得当所有事情处理结束后调用 cancel，告知 durationCtx 的后续 Context 结束
	//defer cancel()

	// 这个 channal 负责通知结束
	finish := make(chan struct{}, 1)
	// 这个 channel 负责通知 panic 异常
	panicChan := make(chan interface{}, 1)

	go func() {
		// 这里增加异常处理
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// 这里做具体的业务
		time.Sleep(10 * time.Second)
		//c.Json(200, "ok")
		// 新的 goroutine 结束的时候通过一个 finish 通道告知父 goroutine
		finish <- struct{}{}
	}()

	select { // 监听 panic
	//case p := <-panicChan:
	// c.Json(500, "panic")
	//监听结束事件
	case <-finish:
		fmt.Println("finish")
		// 监听超时事件
		//case <-durationCtx.Done():
		// c.Json(500, "time out")
	}
}
