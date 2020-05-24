package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
	"time"
)

//logtaill模块采集数据之后发到在这个chanel
type logData struct {
	topic string
	logInfor string
}

var (
	client sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addres []string,chanMax int) (err error) {
	//初始化producer配置对象
	config := sarama.NewConfig()
	//选择leader partion的方式使用轮询
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//producer发送消息成功后，需要leader和follower全部回应ACK
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	//连接kafka集群中的的broker
	client, err = sarama.NewSyncProducer(addres, config)
	if err != nil {
		fmt.Println("连接kafka失败", err)
		return
	}
	fmt.Println("Access kafka successfully")
	//持续收集日志无需关闭！
	//defer client.Close()

	//初始化全局的logDataChan,发往kafka
	logDataChan=make(chan *logData,chanMax)
	//开启后台的gorutine等待logDataChan中数据到来， 然后发往kafka
	go getChanSendToKafka()
	return

}

//给外部暴露1个函数，该函数只把日志数据发送到1个channel中，并非kafka
func SendToChan(topic,line string){
	logData:=&logData{
		topic: topic,
		logInfor:line,
	}
	logDataChan <-logData
}

//从channel获取消息 往kafaka中发
func getChanSendToKafka(){
	for {
		select {
			case logInfo:=<-logDataChan:
				msg:=&sarama.ProducerMessage{}
				msg.Topic=logInfo.topic
				msg.Value=sarama.StringEncoder(logInfo.logInfor)
				pid,offset,err:=client.SendMessage(msg)
				if err!=nil{
					fmt.Println(msg,"消息发送失败",err)
					return
				}
				fmt.Printf("分区ID:%v, offset:%v \n", pid, offset)
				fmt.Println("消息发送成功！")

		default:
			time.Sleep(time.Microsecond*30)

		}
		
	}

}
