package main

import (
	"github.com/hashicorp/consul/api"
)

func main() {
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	consul.Agent().ServiceDeregister("b731c227-d363-4597-b2db-380dfebfdce1")
}
