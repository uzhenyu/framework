package consul

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"github.com/uzhenyu/framework/config"
	"net"
)

func GetIp() (ip []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, addr := range addrs {
		ipNet, isVailIpNet := addr.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = append(ip, ipNet.IP.String())
			}
		}

	}
	return ip
}

// 注册
func NewClient(port int64, address, name, fileName string) error {
	err := config.ReadConfig(fileName)
	if err != nil {
		return err
	}
	c, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", viper.GetString("Nacos.Ip"), "8500")})
	if err != nil {
		return err
	}
	e := GetIp()
	address = e[0]
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
func NewClients(name, fileName string) (string, int64, error) {
	err := config.ReadConfig(fileName)
	if err != nil {
		return "", 0, err
	}
	c, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", viper.GetString("Nacos.Ip"), "8500")})
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
