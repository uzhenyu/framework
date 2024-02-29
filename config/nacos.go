package config

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

const ip = "127.0.0.1"
const port = 8848

var client config_client.IConfigClient
var globalConfig string

func GetClient() error {
	var err error

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithContextPath("/nacos")),
	}

	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(""),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)

	client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return err
}

func GetConfig(group, dataID string) (string, error) {
	if client == nil {
		err := GetClient()
		if err != nil {
			return "", err
		}
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		return "", err
	}

	return content, nil
}

func ListenConfig(group, dataID string) (error, string) {
	if client == nil {
		err := GetClient()
		if err != nil {
			return err, getGlobalConfig()
		}
	}
	err := client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("Config changed, group: " + group + ", dataId: " + dataId + ", content: " + data)

			globalConfig = data
			logs.Info("Updated global config: " + globalConfig)

			// 在这里可以加入配置变化后的逻辑处理
		},
	})
	if err != nil {
		return err, getGlobalConfig()
	}

	return nil, getGlobalConfig()
}

func getGlobalConfig() string {
	if globalConfig == "" {
		return ""
	}
	return globalConfig
}
