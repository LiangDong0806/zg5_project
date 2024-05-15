package global

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	server "zg5/work01/server/proto"
)

var (
	Engine       *gin.Engine
	NacosConf    *NacosConfig
	ApiConf      *ApiConfig
	ConsulClient *api.Client
	ServerClient server.ServerClient
)

type ApiConfig struct {
	Mysql   *MysqlConfig  `yaml:"mysql"`
	Consul  *ConsulConfig `yaml:"consul"`
	ApiConf *GrpcConfig   `yaml:"apiconf"`
}
type NacosConfig struct {
	NamespaceId string `yaml:"NamespaceId"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
}

type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Library  string `yaml:"dbname"`
}

type ConsulConfig struct {
	Id             string   `yaml:"id"`
	Name           string   `yaml:"name"`
	Tags           []string `yaml:"tags"`
	UserConsulName string   `yaml:"userconsulname"`
}

type GrpcConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
