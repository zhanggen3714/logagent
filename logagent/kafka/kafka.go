package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
)

var (
	client sarama.SyncProducer
)

func Init(addres []string) (err error) {
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
	//持续收集日志无需关闭！
	//defer client.Close()
	return

}

func SendToKafka(topic,data string){
	msg:=&sarama.ProducerMessage{}
	msg.Topic=topic
	msg.Value=sarama.StringEncoder(data)
	pid,offset,err:=client.SendMessage(msg)
	if err!=nil{
		fmt.Println(msg,"消息发送失败",err)
		return
	}
	fmt.Printf("分区ID:%v, offset:%v \n", pid, offset)
	fmt.Println("消息发送成功！")
}
