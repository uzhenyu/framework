package grpc

import (
	"encoding/json"
	"fmt"
	"github.com/uzhenyu/framework/consul"
	"github.com/uzhenyu/framework/mysql"
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

func getConfig(serviceName string) (*T, error) {
	configInfo, err := mysql.GetConfig("DEFAULT_GROUP", serviceName)
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

func GetGrpc(serviceName string, register func(s *grpc.Server)) error {
	mysql.Services("127.0.0.1", 8081)
	cof, err := getConfig(serviceName)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", cof.App.Ip, cof.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()
	//反射接口支持查询
	reflection.Register(s)
	consul.NewClient(8081, "10.2.171.80", "wzy")
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return err
}
