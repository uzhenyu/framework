package mysql

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/uzhenyu/framework/config"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type mysqlConfig struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

func Listen(serviceName string) *mysqlConfig {
	globalConfig, _ := config.ListenConfig("DEFAULT_GROUP", serviceName)
	logs.Info(globalConfig, 110)
	if globalConfig != "" {
		data := &mysqlConfig{}
		yamlData := []byte(globalConfig)
		err := yaml.Unmarshal(yamlData, data)
		if err != nil {
			panic(err)
		}
		return data
	}
	return nil
}

func InitMysql(serviceName string) error {
	mysqlCfl := Listen(serviceName)
	logs.Info(mysqlCfl, 0101010101010101010)
	if mysqlCfl != nil {
		type Val struct {
			Mysql mysqlConfig `yaml:"Mysql"`
		}
		mysqlConfigVal := Val{}
		content, err := config.GetConfig("DEFAULT_GROUP", serviceName)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal([]byte(content), &mysqlConfigVal)
		if err != nil {
			fmt.Println("Error: failed to unmarshal YAML content")
			return err
		}
		fmt.Println(content)
		fmt.Println(mysqlConfigVal)
		configM := mysqlConfigVal.Mysql
		dsn := fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
			configM.Username,
			configM.Password,
			configM.Host,
			configM.Port,
			configM.Database,
		)

		// 关闭原有的数据库连接
		if DB != nil {
			sqlDB, _ := DB.DB()
			sqlDB.Close()
		}

		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		logs.Info(dsn, 111111111111111)
		return err
	} else {
		return fmt.Errorf("mysqlCfl id nil")
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
