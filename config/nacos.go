package config

//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/beego/beego/v2/core/logs"
//	"github.com/nacos-group/nacos-sdk-go/v2/clients"
//	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
//	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
//	"github.com/nacos-group/nacos-sdk-go/v2/vo"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//const ip = "127.0.0.1"
//const port = 8848
//
//var DB *gorm.DB
//var client config_client.IConfigClient
//
//type Nachos struct {
//	Username string `json:"Username"`
//	Password string `json:"Password"`
//	Host     string `json:"Host"`
//	Port     string `json:"Port"`
//	Database string `json:"Database"`
//}
//
//func GetClient() error {
//	var err error
//
//	sc := []constant.ServerConfig{
//		*constant.NewServerConfig(ip, port, constant.WithContextPath("/nacos")),
//	}
//
//	//create ClientConfig
//	cc := *constant.NewClientConfig(
//		constant.WithNamespaceId(""),
//		constant.WithTimeoutMs(5000),
//		constant.WithNotLoadCacheAtStart(true),
//		constant.WithLogDir("/tmp/nacos/log"),
//		constant.WithCacheDir("/tmp/nacos/cache"),
//		constant.WithLogLevel("debug"),
//	)
//
//	client, err = clients.NewConfigClient(
//		vo.NacosClientParam{
//			ClientConfig:  &cc,
//			ServerConfigs: sc,
//		},
//	)
//	return err
//}
//
//func GetConfig(group, dataID string) (string, error) {
//	content, err := client.GetConfig(vo.ConfigParam{
//		DataId: dataID,
//		Group:  group,
//	})
//	if err != nil {
//		return "", err
//	}
//
//	return content, nil
//}
//
//// TODO:完成对mysql的监听
//func ListenConfig(group, dataID string) error {
//	return client.ListenConfig(vo.ConfigParam{
//		DataId: dataID,
//		Group:  group,
//		OnChange: func(namespace, group, dataId, data string) {
//			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
//			var nacos Nachos
//			err := json.Unmarshal([]byte(data), &nacos)
//			if err != nil {
//				return
//			}
//			dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
//				nacos.Username,
//				nacos.Password,
//				nacos.Host,
//				nacos.Port,
//				nacos.Database,
//			)
//			UpdateDb(dsn)
//			logs.Info(nacos)
//		},
//	})
//}
//
//func UpdateDb(config string) {
//	Dbs, _ := DB.DB()
//	if Dbs != nil {
//		Dbs.Close()
//	}
//	var err error
//	DB, err = gorm.Open(mysql.Open(config), &gorm.Config{})
//	if err != nil {
//		return
//	}
//}
