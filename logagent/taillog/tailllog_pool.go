package taillog

import (
	"fmt"
	"time"
	"jd.com/logagent/etcd"
)

//定义1个全局的taskpool
var tiallpoolObj *tailPool

type tailPool struct {
	//保存从etcd中获取的所有logAgent配置
	logConfigs            []*etcd.LogEntry
	taskMaping            map[string]*TaillTask
	watchNewConfigChannel chan []*etcd.LogEntry
}

func InitTaskPool(logConfigList []*etcd.LogEntry) {
	//初始化1个taill连接池
	tiallpoolObj = &tailPool{
		logConfigs:            logConfigList, //把当前获取的配置项保存起来
		taskMaping:            make(map[string]*TaillTask,32),
		watchNewConfigChannel: make(chan []*etcd.LogEntry), //无缓冲区通道（没有值1一直阻塞）

	}
	//在pool中初始化日志采集task
	for _, cfg := range logConfigList {
		//生成真正的日志采集模块
		var taiiobj TaillTask
		task,err:=taiiobj.NewTaillTask(cfg.Path, cfg.Topic)
		taskKey := fmt.Sprint(cfg.Path, cfg.Topic)
		tiallpoolObj.taskMaping[taskKey] = task
		if err!=nil {
			fmt.Printf("初始化%s采集日志模块失败%s",cfg.Path,err)
		}

	}
	//开启日志采集池
	go tiallpoolObj.run()

}

//watchNewConfigChannel 配置更新之后，做对应的处理
//1.配置新增
//2.配置删除
//3.配置变更

func (T *tailPool)seekDifference(confs[]*etcd.LogEntry,confsMap map[string]*TaillTask)(difference []*etcd.LogEntry){
	for _, conf := range confs{
		MK := fmt.Sprint(conf.Path, conf.Topic)
		_, ok := confsMap[MK]
		if !ok {
			fmt.Println("检查到key发生变化----->：",MK)
			difference= append(difference,conf)
		}
		continue

	}
	return
}

func (T *tailPool) run() {
	for {
		select {
		case newConfSlice := <-T.watchNewConfigChannel:
			fmt.Println("---------新的配置来了---------")
			//增加
			addList:=T.seekDifference(newConfSlice,T.taskMaping)
			for _, conf := range addList {
				var taiiobj TaillTask
				task,err:=taiiobj.NewTaillTask(conf.Path, conf.Topic)
				if err!=nil {
					fmt.Println("初始化采集日志模块失败",err)
				}
				taskKey := fmt.Sprint(conf.Path, conf.Topic)
				tiallpoolObj.taskMaping[taskKey] = task
				fmt.Printf("增加%s日志采集模块成功\n",taskKey)
			}
			//删除
			newTaskMaping:=make(map[string]*TaillTask,32)
			for _, cnf := range newConfSlice {
				mk:=fmt.Sprint(cnf.Path,cnf.Topic)
				newTaskMaping[mk]=nil
			}
			deleteList:=T.seekDifference(T.logConfigs,newTaskMaping)
			for _, item := range deleteList {
				taskKey := fmt.Sprint(item.Path, item.Topic)
				T.taskMaping[taskKey].exit()
				delete(T.taskMaping,taskKey)
			}
			//更新logConfigs
			tiallpoolObj.logConfigs=newConfSlice
			fmt.Println("最新日志采集任务列表",T.taskMaping)
		default:
			time.Sleep(time.Second)

		}

	}
}

//向外暴露watchNewConfigChannel
func PushNewConfig() chan<- []*etcd.LogEntry {
	return tiallpoolObj.watchNewConfigChannel
}
