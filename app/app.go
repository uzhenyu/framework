package app

import (
	"github.com/uzhenyu/framework/config"
	"github.com/uzhenyu/framework/mysql"
)

func Init(serviceName, fileName string, apps ...string) error {
	var err error
	err = config.GetClient(fileName)
	if err != nil {
		return err
	}
	for _, val := range apps {
		switch val {
		case "mysql":
			err = mysql.InitMysql(serviceName, fileName)
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}
