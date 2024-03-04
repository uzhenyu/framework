package app

import (
	"github.com/uzhenyu/framework/mysql"
)

func Init(serviceName string, apps ...string) error {
	var err error
	err = mysql.GetClient()
	if err != nil {
		return err
	}
	for _, val := range apps {
		switch val {
		case "mysql":
			err = mysql.InitMysql(serviceName)
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}
