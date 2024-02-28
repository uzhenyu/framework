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

func GetClient() error {
	var err error

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, port, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
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

var globalConfig string

// TODO:完成对mysql的监听
func ListenConfig(group, dataID string) (error, string) {
	if client == nil {
		err := GetClient()
		if err != nil {
			return err, ""
		}
	}
	err := client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			// 在配置发生变化时，更新配置
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)

			// 更新全局变量保存最新配置
			globalConfig = data
			logs.Info(globalConfig, 222222222222222)
			// 在这里可以添加你的配置更新逻辑，例如更新数据库连接、重新加载配置文件等
			// updateConfig(data)
		},
	})
	if err != nil {
		return err, globalConfig
	}
	return nil, globalConfig
}
