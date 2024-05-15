package global

import "github.com/hashicorp/consul/api"

var (
	RpcALLConf   *RpcAllClient
	NacosConf    *NacosConfig
	ConsulClient *api.Client
)

type RpcAllClient struct {
	Grpc   *GrpcConfig   `yaml:"grpc"`
	Consul *ConsulConfig `yaml:"consul"`
	Mysql  *MysqlConfig  `yaml:"mysql"`
}

type NacosConfig struct {
	NamespaceId string `yaml:"NamespaceId"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
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
