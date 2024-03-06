package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const ip = "10.2.171.13"
const port = 8848

var DB *gorm.DB
var client config_client.IConfigClient

type Nachos struct {
	Mysql struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Database string `json:"Database"`
	} `json:"Mysql"`
}

func InitMysql(serviceName string) error {
	err := GetClient()
	if err != nil {
		return err
	}

	config, err := GetConfig(serviceName, "wzy")
	if err != nil {
		return err
	}
	var nacos Nachos
	err = json.Unmarshal([]byte(config), &nacos)
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		nacos.Mysql.Username,
		nacos.Mysql.Password,
		nacos.Mysql.Host,
		nacos.Mysql.Port,
		nacos.Mysql.Database,
	)

	err = ListenConfig(serviceName, "wzy")
	if err != nil {
		return err
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

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
	if err != nil {
		return err
	}
	return err
}

func GetConfig(group, dataID string) (string, error) {
	if client == nil {
		// 初始化 client
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

func ListenConfig(group, dataID string) error {
	return client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			var nacos Nachos
			err := json.Unmarshal([]byte(data), &nacos)
			if err != nil {
				return
			}
			dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
				nacos.Mysql.Username,
				nacos.Mysql.Password,
				nacos.Mysql.Host,
				nacos.Mysql.Port,
				nacos.Mysql.Database,
			)
			UpdateDb(dsn)
			logs.Info(nacos)
		},
	})
}

func UpdateDb(config string) {
	Dbs, _ := DB.DB()
	if Dbs != nil {
		Dbs.Close()
	}
	var err error
	DB, err = gorm.Open(mysql.Open(config), &gorm.Config{})
	if err != nil {
		return
	}

}

func WithTX(txFc func(tx *gorm.DB) error) {
	var err error
	tx := DB.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}

func Services(ips string, ports int64) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   port,
		},
	}
	cli, _ := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	_, err := cli.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ips,
		Port:        uint64(ports),
		ServiceName: "wzy",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "DEFAULT",       // 默认值DEFAULT
		GroupName:   "DEFAULT_GROUP", // 默认值DEFAULT_GROUP
	})
	if err != nil {
		return
	}
}
