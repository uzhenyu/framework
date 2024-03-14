package grpc

import (
	"google.golang.org/grpc"
)

func Client(toService, fileName string) (*grpc.ClientConn, error) {
	//cnfStr, err := config.GetConfig("DEFAULT_GROUP", toService)
	//if err != nil {
	//	return nil, err
	//}
	//cnf := new(T)
	//err = json.Unmarshal([]byte(cnfStr), &cnf)
	//if err != nil {
	//	return nil, err
	//}
	//logs.Info(cnf.App.Ip, cnf.App.Port)
	//err := config.ReadConfig(fileName)
	//address, port, err := consul.NewClients(viper.GetString("Wzy.DataID"), fileName)
	//if err != nil {
	//	return nil, err
	//}
	return grpc.Dial("consul://10.2.171.70:8500/"+"wzy"+"?wait=14s", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))
}
