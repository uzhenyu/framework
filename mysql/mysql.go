package mysql

import (
	"encoding/json"
	"fmt"
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

func InitMysql(serviceName string) error {
	err := config.GetClient()
	if err != nil {
		return err
	}

	configs, err := config.GetConfig(serviceName, "wzy")
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

	err = config.ListenConfig(serviceName, "wzy")
	if err != nil {
		return err
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	defer func() {
		db, err := DB.DB()
		if err != nil {
			return
		}
		_ = db.Close()
	}()
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
