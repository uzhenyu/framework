package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/spf13/viper"
	"github.com/uzhenyu/framework/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Nachos struct {
	Mysql struct {
		Username string `json:"Username"`
		Password string `json:"Password"`
		Host     string `json:"Host"`
		Port     string `json:"Port"`
		Database string `json:"Database"`
	} `json:"Mysql"`
}

func InitMysql(serviceName, fileName string) error {
	err := config.GetClient(fileName)
	if err != nil {
		return err
	}
	err = config.ReadConfig(fileName)
	if err != nil {
		return err
	}
	configs, err := config.GetConfig(serviceName, viper.GetString("Wzy.DataID"), fileName)
	if err != nil {
		return err
	}
	var nacos Nachos
	err = json.Unmarshal([]byte(configs), &nacos)
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
	logs.Info(viper.GetString("Wzy.DataID"))
	err = config.ListenConfig(serviceName, viper.GetString("Wzy.DataID"))
	if err != nil {
		return err
	}
	logs.Info(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	//defer func() {
	//	db, err := DB.DB()
	//	if err != nil {
	//		return
	//	}
	//	_ = db.Close()
	//}()
	return nil
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
