package consule

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"zg5/work01/server/common/global"
)

type RegistryConsul struct {
	Host string
	Port int
}

type RegistryClient interface {
	RegisterConsul(id string, name string, tags []string) error
	FilterConsulByName(name string) (map[string]*api.AgentService, error)
	DeregisterConsul(id string) error
}

func NewConsulClient(host string, port int) RegistryClient {
	return &RegistryConsul{
		Host: host,
		Port: port,
	}
}

// 注册consul
func (r *RegistryConsul) RegisterConsul(id string, name string, tags []string) error {
	err := global.ConsulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    global.ServerConfig.Grpc.Port,
		Address: global.ServerConfig.Grpc.Host,
	})
	return err
}

// 根据服务名称过滤
func (r *RegistryConsul) FilterConsulByName(name string) (map[string]*api.AgentService, error) {
	log.Println(fmt.Sprintf(`"Service == %s"`, name), "[[[")
	return global.ConsulClient.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, "Server"))
}

// 注销服务
func (r *RegistryConsul) DeregisterConsul(id string) error {
	return global.ConsulClient.Agent().ServiceDeregister(id)
}
