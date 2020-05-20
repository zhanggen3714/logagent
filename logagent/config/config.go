package config

type APPConf struct {
	KafkaConf `ini:"kafka""`
	TaillogConfig `ini:"taillog""`
}

type KafkaConf struct {
	Address string `ini:"address""`
	Topic string   `ini:"topic""`
}

type TaillogConfig struct {
	Nigix string `ini:"nigix""`
}
