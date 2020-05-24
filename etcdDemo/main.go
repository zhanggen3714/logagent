package main

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"fmt"
	"time"
)

func main() {
	//创建1个连接etcd的连接
	client,err:=clientv3.New(clientv3.Config{
		Endpoints: []string{"192.168.56.133:2371"},
		DialTimeout: 5*time.Second,
	})
	if err!=nil{
		fmt.Println("I try to connect etcd fiald ",err)
	}
	fmt.Println("The connection to etcd was connected!")
	//记得关闭连接
	defer client.Close()

	//put值
	//	{"path":"D:\\tmp\\nginx.txt","topic":"ngix-log"},
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second)
	//设置1秒钟超时时间
	key:="tailconfig"
	//{"path":"D:\\tmp\\redis.txt","topic":"redis-log"},
	jsonString:=`
		[
			{"path":"D:\\tmp\\nginx.txt","topic":"ngix-log"},
			{"path":"D:\\tmp\\mysql.txt","topic":"mysql-log"}
			]
		`
	_,err=client.Put(ctx,key,jsonString)
	cancel()
	if err!=nil{
		fmt.Println("The instraction put value to etcd was faild.")
		return
	}
	////get值
	//ctx,cancel=context.WithTimeout(context.Background(),time.Second)
	////查询所有以name开头的key（支持模糊查询）
	//response,err:=client.Get(ctx,key,clientv3.WithPrefix())
	//cancel()
	//if err!=nil{
	//	fmt.Println("The instraction get value from etcd was faild.")
	//}
	////fmt.Println(response)
	///*
	//&{cluster_id:2654831503084648812 member_id:3952383196236056773 revision:3 raft_term:2  [key:"name" create_revision:2 mod_revision:3 version:2 value:",martin"
	//] false 1 {} [] 0}
	//*/
	//for _,ev:=range response.Kvs{
	//	fmt.Printf("key:%s value:%s\n",ev.Key,ev.Value)
	//}
	//ctx=context.TODO()
	////安插1个探针 监控name key，返回1个通道
	//monitorChanel:=client.Watch(ctx,key)
	////不断从通道中获取监控事件
	//for monitor:= range monitorChanel{
	//	//获取监控事件
	//	for _,event:=range monitor.Events{
	//		fmt.Printf("type:%v,key:%v,value:%v\n",event.Type,string(event.Kv.Key),string(event.Kv.Value))
	//	}
	//
	//}


}
