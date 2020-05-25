package config

//logagent app配置
type APPConf struct {
	KafkaConf `ini:"kafka"`
	EtcdConfig `ini:"etcd"`

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
	Storekey string `ini:"storekey"`
}


