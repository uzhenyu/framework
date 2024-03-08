package consul

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

// 注册
func NewClient(port int64, address, name string) error {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	err = c.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.New().String(),
		Name:    name,
		Tags:    []string{"GRPC"},
		Port:    int(port),
		Address: address,
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",                                //间隔时常
			Timeout:                        "5s",                                //退出
			GRPC:                           fmt.Sprintf("%v:%v", address, port), //
			DeregisterCriticalServiceAfter: "30s",                               //注销
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// 获取健康服务
func NewClients(name string) (string, int64, error) {
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return "", 0, err
	}
	byName, data, err := c.Agent().AgentHealthServiceByName(name)
	if err != nil {
		return "", 0, err
	}
	fmt.Println(byName)
	if byName != "passing" {

	}
	return data[0].Service.Address, int64(data[0].Service.Port), nil
}
