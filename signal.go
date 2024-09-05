package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func DealSignalHup() {
	log.Println("处理HUP信号")
}

func CatchSighup() {
	// 创建一个通道用于接收信号
	signalChan := make(chan os.Signal, 1)

	// 捕捉 SIGHUP 信号
	signal.Notify(signalChan, syscall.SIGHUP)

	// 启动一个协程来监听信号
	go func() {
		for {
			sig := <-signalChan
			log.Println("收到信号:", sig)

			// 处理 SIGHUP 信号，重新加载配置
			if sig == syscall.SIGHUP {
				DealSignalHup()
			}
		}
	}()
}
