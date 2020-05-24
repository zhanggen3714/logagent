package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"jd.com/logagent/config"
	"jd.com/logagent/etcd"
	"jd.com/logagent/kafka"
	"jd.com/logagent/taillog"
	"sync"
	"time"
)

var (
	//new初始化1个指针
	cfg = new(config.APPConf)
)

func main() {

	//0.加载配置文件映射成 struct
	ini.MapTo(cfg, "./config/config.ini")
	//1.初始化kafka连接和 taill-->chanel <-->kafka队列:
	err := kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.ChanSize)
	if err != nil {
		fmt.Println("初始化kafka失败", err)
		return
	}
	//2.初始化1个etcd
	err = etcd.Init(cfg.EtcdConfig.Address, time.Duration(cfg.EtcdConfig.Timeout)*time.Second)
	if err != nil {
		fmt.Println("初始化etcd失败", err)
		return
	}
	//去kafka更新logagent的配置（开发环境使用）
	etcd.TryPut(cfg.TailLog.Storekey)
	logConfigList, err := etcd.Getconfig(cfg.TailLog.Storekey)
	if err != nil {
		fmt.Println("Get config from etcd faild%s\n", err)
	}


	//查看所有logant配置项
	fmt.Println("Get config from etcd successsfully")
	for _,logconfig:= range logConfigList {
			fmt.Println(logconfig.Topic,logconfig.Path)
	}

	//3.从etcd中获取logAgent配置，根据logant配置项，初始1个日志采集池（创建x个日志采集对象taillobj）
	taillog.InitTaskPool(logConfigList)

	//4.后台开启1个哨兵循环监控etcd中存储loAant配置项的key
	newConfChan:=taillog.PushNewConfig()
	var wg sync.WaitGroup
	wg.Add(1)
	//在etcd包中watch etcd服务的非删除操作，通过chnnel通知给的日志采集池（InitTaskPool）动态管理task
	go etcd.Watchconfig(cfg.TailLog.Storekey,newConfChan)
	wg.Wait()

}
