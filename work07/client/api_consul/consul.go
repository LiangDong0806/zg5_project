package api_consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"zg5/work/work07/client/api_global"
)

type RegistryConsul struct {
	Host string
	Port int
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
	err := api_global.ConsulClient.Agent().ServiceRegister(&api.AgentServiceRegistration{
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
	return api_global.ConsulClient.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, name))
}

func (r *RegistryConsul) DeregisterConsulById(id string) error {
	//TODO implement me
	return api_global.ConsulClient.Agent().ServiceDeregister(id)
}

// 服务发现
func (r *RegistryConsul) AgentHealthServiceByName(name string) []api.AgentServiceChecksInfo {
	name, i, err := api_global.ConsulClient.Agent().AgentHealthServiceByName(name)
	if err != nil {
		panic(err)
	}
	if name != "passing" {
		panic("is not health utils")
	}
	return i
}
