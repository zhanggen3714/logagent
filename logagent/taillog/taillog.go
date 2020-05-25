package taillog

import (
	"context"
	"github.com/hpcloud/tail"
	"fmt"
	"jd.com/logagent/kafka"
)

//1个具体的日志收集任务（TaillTask）
type TaillTask struct {
	path     string
	topic    string
	instance *tail.Tail
	//exit task
	ctx  context.Context
	exit context.CancelFunc
}

//实例化1个具体的日志收集任务（TaillTask）
func (T *TaillTask) NewTaillTask(path, topic string)(task *TaillTask,err error){
	ctx, cancel := context.WithCancel(context.Background())
	task=&TaillTask{path: path, topic: topic, ctx: ctx, exit: cancel,}
	//taill 文件配置
	config := tail.Config{
		ReOpen:    true,                                 //重新打开文件
		Follow:    true,                                 //跟随文件
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从文件的哪个地方开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,
	}
	//给task任务填充taill（1个具体打开文件的taillobj）
	task.instance, err = tail.TailFile(task.path, config)
	if err != nil {
		fmt.Println("文件打开失败", err)

	}
	//直接去采集日志
	go task.run()
	return
}

//从tailobj中读取日志内容---->kafka topic方法
func (T *TaillTask)run() {
	fmt.Printf("开始收集%s日志\n",T.path)
	for {
		select {
		//父进程调用了cancel
		case <-T.ctx.Done():
			fmt.Printf("taill任务%s%s退出了...\n",T.topic,T.path)
			return
		case line := <-T.instance.Lines:
			fmt.Printf("从%s文件中获取到内容%s",T.path,line.Text)
			//taill采集到数据-----channel------>kafka 异步
			kafka.SendToChan(T.topic, line.Text)
		}

	}

}
