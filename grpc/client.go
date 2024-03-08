package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/uzhenyu/framework/config"
	"github.com/uzhenyu/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(toService string) (*grpc.ClientConn, error) {
	cnfStr, err := config.GetConfig("DEFAULT_GROUP", toService)
	if err != nil {
		return nil, err
	}
	cnf := new(T)
	err = json.Unmarshal([]byte(cnfStr), &cnf)
	if err != nil {
		return nil, err
	}
	logs.Info(cnf.App.Ip, cnf.App.Port)
	return grpc.Dial(fmt.Sprintf("%v:%v", cnf.App.Ip, cnf.App.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func ConnectionClient() (*grpc.ClientConn, error) {
	address, port, err := consul.NewClients("wzy")
	if err != nil {
		return nil, err
	}
	return grpc.Dial(fmt.Sprintf("%v:%v", address, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
