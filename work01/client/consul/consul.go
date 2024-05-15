package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"zg5/work01/client/global"
)

type RegistryConsul struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RegistryClient interface {
	RegisterConsul(id, name string, tags []string) error
	FilterConsulByName(name string) (map[string]*api.AgentService, error)
	DeregisterConsulById(id string) error
	AgentHealthServiceByName(name string) []api.AgentServiceChecksInfo
}

func NewConsulClient(host string, port int) RegistryClient {
	return &RegistryConsul{
		Host: host,
		Port: port,
	}
}

func (r *RegistryConsul) RegisterConsul(id, name string, tags []string) error {
	//TODO implement me
	err := global.ConsulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    r.Port,
		Address: r.Host,
	})
	return err
}

func (r *RegistryConsul) FilterConsulByName(name string) (map[string]*api.AgentService, error) {
	//TODO implement me
	return global.ConsulClient.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, name))
}

func (r *RegistryConsul) DeregisterConsulById(id string) error {
	//TODO implement me
	return global.ConsulClient.Agent().ServiceDeregister(id)
}

func (r *RegistryConsul) AgentHealthServiceByName(name string) []api.AgentServiceChecksInfo {
	//TODO implement me
	name, b, _ := global.ConsulClient.Agent().AgentHealthServiceByName(name)
	if name != "passing" {
		panic("is not health utils")
	}
	return b
}
