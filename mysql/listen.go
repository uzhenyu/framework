package mysql

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/uzhenyu/framework/config"
	"gopkg.in/yaml.v2"
)

type mysqll struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
}

type mysqls struct {
	Mysql struct {
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		Database string `yaml:"Database"`
	} `yaml:"Mysql"`
}

func Listen(serviceName string) mysqll {
	err, globalConfig := config.ListenConfig("DEFAULT_GROUP", serviceName)
	if globalConfig != "" {
		yamlData := `
			Mysql:
			  Username: root
			  Password: 12345
			  Host: 127.0.0.1
			  Port: 3306
			  Database: zg5 222222222222222
			`
		var config mysqls
		err = yaml.Unmarshal([]byte(yamlData), &config)
		if err != nil {
			logs.Info("解析失败")
		}
		mysqlCfl := mysqll{}
		mysqlCfl.Port = config.Mysql.Port
		mysqlCfl.Username = config.Mysql.Username
		mysqlCfl.Password = config.Mysql.Password
		mysqlCfl.Host = config.Mysql.Host
		mysqlCfl.Database = config.Mysql.Database
		type Val struct {
			Mysql mysqlConfig `yaml:"Mysql"`
		}
		return mysqlCfl
	}
	return mysqll{}
}
