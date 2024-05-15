package api_global

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"zg5/work/work07/server/proto/server"
)

var (
	Router       *gin.Engine
	ClientConfig *ClientConf
	ConsulClient *api.Client
	ServerClient server.ServerClient
	NacosConfig  *NacosConf
)

type ClientConf struct {
	Consul  *ConsulConf  `yaml:"consul"`
	ApiConf *ApiConf     `yaml:"apiconf"`
	Mysql   *MysqlConfig `yaml:"mysql"`
}

type NacosConf struct {
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	NamespaceId string `yaml:"NamespaceId"`
}

type ConsulConf struct {
	Id             string   `yaml:"Id"`
	Name           string   `yaml:"name"`
	Tags           []string `yaml:"tags"`
	UserConsulName string   `yaml:"userconsulname"`
}

type ApiConf struct {
	Host string `json:"Host" yaml:"host"`
	Port int    `json:"Port" yaml:"port"`
}
type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}
