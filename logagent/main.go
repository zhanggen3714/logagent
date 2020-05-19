package main

import (
	"time"
	"fmt"
	"jd.com/logagent/kafka"
	"jd.com/logagent/taillog"
)

func run() {
	//读取日志
	for {
		select {
		case line := <-taillog.ReadChanel():
			//往kafka中发送获取到的日志信息
			kafka.SendToKafka("web-log", line.Text)

		default:
			time.Sleep(time.Second)
		}

	}

}

func main() {
	err := kafka.Init([]string{"192.168.56.133:9092"})
	if err != nil {
		fmt.Println("初始化kafka失败", err)
		return
	}
	err = taillog.Init("D:\\mylog.txt")
	if err != nil {
		fmt.Println("初始化tail失败", err)
		return
	}

	run()

}
