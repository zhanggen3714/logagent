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
func (T *tailPool) run() {
	for {
		select {
		case newConfSlice := <-T.watchNewConfigChannel:
			fmt.Println("新的配置来了！", newConfSlice)
			//发现需要新增的task
			for _, conf := range newConfSlice {
				taskKey := fmt.Sprint(conf.Path, conf.Topic)
				_, ok := T.taskMaping[taskKey]
				if ok {
					//原来有
					continue
				} else {
					//新增的 task
					var taiiobj TaillTask
					task,err:=taiiobj.NewTaillTask(conf.Path, conf.Topic)
					if err!=nil {
						fmt.Println("初始化采集日志模块失败",err)
					}
					taskKey := fmt.Sprint(conf.Path, conf.Topic)
					tiallpoolObj.taskMaping[taskKey] = task
					fmt.Printf("新增%s日志采集模块成功",task.path)
				}
			}

			//发现删除/停止task
			for _, oldConf := range T.logConfigs {
				isDelete := true
				for _, newConf := range newConfSlice {
					if oldConf.Path == newConf.Path && oldConf.Topic == oldConf.Topic {
						isDelete = false
						continue
					}
				}
				if isDelete{
					taskKey := fmt.Sprint(oldConf.Path, oldConf.Topic)
					T.taskMaping[taskKey].exit()
					delete(T.taskMaping,taskKey)
				}
			}
			//更新logConfigs
			//tiallpoolObj.logConfigs=newConfSlice
		default:
			time.Sleep(time.Second)

		}

	}
}

//向外暴露watchNewConfigChannel
func PushNewConfig() chan<- []*etcd.LogEntry {
	return tiallpoolObj.watchNewConfigChannel
}
