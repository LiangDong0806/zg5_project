package global

import (
	"github.com/hashicorp/consul/api"
)

var (
	ServerConfig *ServerConf //TODO 111
	NacosConfig  *NacosConf  //TODO 222
	ConsulClient *api.Client
)

type ServerConf struct { //TODO 111
	Consul *ConsulConfig `yaml:"consul"`
	Grpc   *GrpcConfig   `yaml:"grpc"`
	Mysql  *MysqlConfig  `yaml:"mysql"`
}

type NacosConf struct { //TODO 222
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	NamespaceId string `yaml:"NamespaceId"`
}

type GrpcConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConsulConfig struct {
	Id   string   `yaml:"id"`
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}

type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}
