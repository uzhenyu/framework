package grpc

import (
	"fmt"
	"github.com/uzhenyu/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(toService string) (*grpc.ClientConn, error) {
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
	address, port, err := consul.NewClients("fzjh")
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", address, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
