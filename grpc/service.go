package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/uzhenyu/framework/config"
	"github.com/uzhenyu/framework/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type T struct {
	App struct {
		Ip     string `json:"ip"`
		Port   string `json:"port"`
		Secret string `json:"secret"`
	} `json:"app"`
}

func getConfig(serviceName, fileName string) (*T, error) {
	configInfo, err := config.GetConfig("DEFAULT_GROUP", serviceName, fileName)
	if err != nil {
		return nil, err
	}
	cnf := new(T)
	err = json.Unmarshal([]byte(configInfo), &cnf)
	if err != nil {
		return nil, err
	}
	return cnf, nil
}

func GetGrpc(serviceName, fileName string, register func(s *grpc.Server)) error {
	//mysql.Services("10.2.171.13", 8081)
	err := config.ReadConfig(fileName)
	if err != nil {
		return err
	}
	cof, err := getConfig(serviceName, fileName)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "0.0.0.0", cof.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	err = consul.NewClient(8081, viper.GetString("Nacos.Ip"), viper.GetString("Wzy.wzy"), fileName)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	//反射接口支持查询
	reflection.Register(s)
	//健康检查
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return err
}
