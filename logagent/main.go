package main

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"jd.com/logagent/config"
	"jd.com/logagent/etcd"
	"jd.com/logagent/kafka"
	"jd.com/logagent/seek_local_ip"
	"jd.com/logagent/taillog"
	"sync"
	"time"
)

var (
	//new初始化1个指针
	cfg = new(config.APPConf)
)

func main() {
	ipString,err:=seek_local_ip.SeekLocalIP()
	if err!=nil{
		fmt.Println("请检查本地网络是否正常")
		return
	}
	//0.加载配置文件映射成 struct
	ini.MapTo(cfg, "./config/config.ini")
	//初始etcd中config key的名称:引入 ip/自定义string 是为了分布式logagent根据不同业务线，从etcd加载配置文件信息.
	configsLine:= flag.String("etcdkey", ipString,"Please type in a string as a part as your key getting config from etcd")
	flag.Parse()
	//1.初始化kafka连接和 taill-->chanel <-->kafka队列:
	err= kafka.Init([]string{cfg.KafkaConf.Address},cfg.KafkaConf.ChanSize)
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
	etcdKey:=fmt.Sprintf(cfg.EtcdConfig.Storekey,*configsLine)
	//去kafka更新logagent的配置（开发环境使用）
	etcd.TryPut(etcdKey)
	//从etcd中拉取loagent的所有配置
	logConfigList,err := etcd.Getconfig(etcdKey)
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
	go etcd.Watchconfig(etcdKey,newConfChan)
	wg.Wait()

}
