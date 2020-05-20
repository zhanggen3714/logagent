package main

import (
	"gopkg.in/ini.v1"
	"jd.com/logagent/config"
	"time"
	"fmt"
	"jd.com/logagent/kafka"
	"jd.com/logagent/taillog"
)

var (
	//new初始化1个指针
	cfg=new(config.APPConf)
)


func run() {
	//读取日志
	for {
		select {
		case line := <-taillog.ReadChanel():
			//往kafka中发送获取到的日志信息
			kafka.SendToKafka(cfg.KafkaConf.Topic, line.Text)

		default:
			time.Sleep(time.Second)
		}

	}

}

func main() {

	//cfg,err:=ini.Load("./config/config.ini")
	//kafkaAdrr:=cfg.Section("kafka").Key("address").String()
	//topic:=cfg.Section("kafka").Key("topic").String()
	//logPath:=cfg.Section("logtaill").Key("path").String()
	//配置文件映射成 struct
	ini.MapTo(cfg,"./config/config.ini")
	//初始化kafka:
	//ps：如果没有struct中嵌套struct的字段没有重名可以简写：cfg.Address
	err:= kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Println("初始化kafka失败", err)
		return
	}
	//初始化1个日志文件
	err = taillog.Init(cfg.TaillogConfig.Nigix)
	if err != nil {
		fmt.Println("初始化tail失败", err)
		return
	}

	run()

}
