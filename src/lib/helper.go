package lib

import (
	"os"
	"os/signal"
	"syscall"
)

func Maintain() {
	// 创建一个传递信号的 channel
	sigChan := make(chan os.Signal, 1)
	// 监听退出信号
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	// 等待信号
	<-sigChan
	// 关闭主线程
	os.Exit(0)
}
