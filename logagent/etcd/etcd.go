package etcd

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"fmt"
	"time"
)

var etcdClient *clientv3.Client

//Unmarshal kafka中的logAgent配置
type LogEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

//初始1个全局的etcd连接
func Init(addr string, timeout time.Duration) (err error) {
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout,
	})
	if err != nil {
		fmt.Println("Keep tring to log on etcd,and keep getting error message like ", err)
	}
	fmt.Println("Access etcd successfully")
	return
}

//根据key从etcd中获取配置项
func Getconfig(key string) (logConfigList []*LogEntry, err error) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	//etcdResponse belog to RangeResponse struct
	etcdResponse, err := etcdClient.Get(cxt, key)
	cancel()
	if err != nil {
		fmt.Printf("get %s faild\n", key)
		return
	}
	for _, ev := range etcdResponse.Kvs {
		err = json.Unmarshal(ev.Value, &logConfigList)
		if err != nil {
			fmt.Printf("json Unmarshal logconfig list faild\n", err)
			return
		}

	}
	return
}

//从main函数中开启Watchconfig协程，在后台一直监控kafka中的key，通知taillpool
func Watchconfig(key string, newConfCh chan<- []*LogEntry) {
	ch := etcdClient.Watch(context.Background(), key)
	for resp := range ch {
		for _, event := range resp.Events {
			fmt.Println(event.Type, string(event.Kv.Key), string(event.Kv.Value))
			var newConf []*LogEntry
			//监听到非删除事件
			if event.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(event.Kv.Value, &newConf)
				if err != nil {
					fmt.Println("json unmarshal faild", err)
					continue
				}
			}
			fmt.Println("Get new config", newConf)
			newConfCh <- newConf

		}
	}

}

//测试更新LogAgent配置，测试环境使用
func TryPut(key string) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	/*
	`

	*/
	jsonString :=
		`
		[
			{"path":"D:\\tmp\\nginx.txt","topic":"ngix-log"},
			{"path":"D:\\tmp\\redis.txt","topic":"redis-log"},
			{"path":"D:\\tmp\\mysql.txt","topic":"mysql-log"}
			]
		`
	_, err := etcdClient.Put(cxt, key, jsonString)
	cancel()
	if err != nil {
		fmt.Printf("put %s faild %s\n", key, err)
	}
	fmt.Printf("put %s key successfully\n", key)
}
