package global

import (
	"github.com/hashicorp/consul/api"
	"github.com/olivere/elastic/v7"
)

var (
	ServerConfig *ServerConf //TODO 111
	NacosConfig  *NacosConf  //TODO 222
	ConsulClient *api.Client
	Client       *elastic.Client
)

type ServerConf struct { //TODO 111
	Consul        *ConsulConfig        `yaml:"consul"`
	MySQL         *MySQLConfig         `yaml:"mysql"`
	Grpc          *GrpcConfig          `yaml:"grpc"`
	ElasticSearch *ElasticsearchConfig `yaml:"elasticsearch"`
}

type NacosConf struct { //TODO 222
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	NamespaceId string `yaml:"NamespaceId"`
}

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Library  string `yaml:"library"`
}
type GrpcConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type ConsulConfig struct {
	Name string   `yaml:"name"`
	Id   string   `yaml:"id"`
	Tags []string `yaml:"tags"`
}
type ElasticsearchConfig struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Index string `yaml:"index"`
}
