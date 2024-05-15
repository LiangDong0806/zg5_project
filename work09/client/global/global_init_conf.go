package global

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"zg5/work/work09/server/proto/server"
)

var (
	ApiALLConf   *ApiAllClient
	NacosConf    *NacosConfig
	ConsulClient *api.Client
	ServerClient server.ServerClient
	Routers      *gin.Engine
)

type ApiAllClient struct {
	ApiConf *ApiConfig    `yaml:"apiconf"`
	Consul  *ConsulConfig `yaml:"consul"`
	Mysql   *MysqlConfig  `yaml:"mysql"`
}

type NacosConfig struct {
	NamespaceId string `yaml:"NamespaceId"`
	DataId      string `yaml:"DataId"`
	Group       string `yaml:"Group"`
	Host        string `yaml:"Host"`
	Port        int    `yaml:"Port"`
}

type ApiConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConsulConfig struct {
	Id             string   `yaml:"id"`
	Name           string   `yaml:"name"`
	Tags           []string `yaml:"tags"`
	UserConsulName string   `yaml:"userconsulname"`
}

type MysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Dbname   string `yaml:"dbname"`
}
