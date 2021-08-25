package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf Config

type Config struct {
	Server    server    `yaml:"server,omitempty"`
	Databases DB `yaml:"databases,omitempty"`
}

type server struct {
	Name string `yaml:"name,omitempty"`
	Port string `yaml:"port,omitempty"`
	Host string `yaml:"host,omitempty"`
}

type DB struct {
	Host     string `yaml:"host,omitempty"`
	Port     int `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Database string `yaml:"database,omitempty"`
	PoolSize int    `yaml:"poolSize,omitempty"`
	Slow     int    `yaml:"slow,omitempty"`
}

func InitConfig() {
	//读取配置
	yamlFile, err := ioutil.ReadFile("conf/config.yaml")
	if err == nil {
		err = yaml.Unmarshal(yamlFile, &Conf)
	} else {
		fmt.Println(err)
	}
	//日志地址
	//file := "./storage/logs/" + Conf.App.Name + "_" + time.Now().Format("20060102") + ".log"
	//logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	//loger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	////SetFlags设置输出选项
	//loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//
	//json.Marshal(err)

}
