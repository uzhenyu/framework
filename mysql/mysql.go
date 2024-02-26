package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() error {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/zg5?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
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
