package mysql

import (
	"fmt"
	"github.com/uzhenyu/framework/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var err error

func InitMysql(serviceName string) error {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.Nachos{}.Username,
		config.Nachos{}.Password,
		config.Nachos{}.Host,
		config.Nachos{}.Port,
		config.Nachos{}.Database,
	)
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func WithTX(txFc func(tx *gorm.DB) error) {
	var err error
	tx := config.DB.Begin()
	err = txFc(tx)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
}
