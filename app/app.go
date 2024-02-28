package app

import "github.com/uzhenyu/framework/mysql"

func Init(apps ...string) error {
	var err error
	for _, val := range apps {
		switch val {
		case "mysql":
			err = mysql.InitMysql("wzy")
		}
	}
	return err
}
