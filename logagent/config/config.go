package config

//logagent app配置
type APPConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConfig `ini:"etcd"`
	TailLog `ini:"taillog"`
}

//kafka配置
type KafkaConf struct {
	Address string `ini:"address"`
	ChanSize int `ini:"log_chan_max_size"`
}
//etcd配置
type EtcdConfig struct {
	Address string `ini:"address"`
	Timeout int   `ini:"timeout"`
}

//logAgeny配置
type TailLog struct {
	Storekey string `ini:"storekey"`
}
