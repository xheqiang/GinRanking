package config

import (
	"fmt"
	"ginRanking/util/logger"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

/* const (
	MYSQL_HOST     = "127.0.0.1"
	MYSQL_PORT     = "3306"
	MYSQL_USER     = "root"
	MYSQL_PASSWORD = "root"
	MYSQL_DB       = "ranking"
	MYSQL_CHARSET  = "utf8mb4"

	MYSQLDB = MYSQL_USER + ":" + MYSQL_PASSWORD + "@tcp(" + MYSQL_HOST + ":" + MYSQL_PORT + ")/" + MYSQL_DB + "?charset=" + MYSQL_CHARSET + "&parseTime=True&loc=Local"

	//MYSQLDB = "root:root@tcp(127.0.0.1:3306)/ranking?charset=utf8mb4&parseTime=True&loc=Local"

	REDIS_HOST     = "127.0.0.1"
	REDIS_PORT     = "6379"
	REDIS_PASSWORD = ""
	REDIS_DB       = 0
	REDIS_ADDR     = REDIS_HOST + ":" + REDIS_PORT
) */

type AppConfig struct {
	MySQLConfig `mapstructure:"mysql"`
	RedisConfig `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	Local        string `mapstructure:"local"`
	MaxIdleTime  int    `mapstructure:"max_idle_time"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

var AppConf = AppConfig{}

func init() {
	workDir, _ := os.Getwd()
	viperConfig := viper.New()
	viperConfig.SetConfigFile(path.Join(workDir, "config/config.yaml"))

	if err := viperConfig.ReadInConfig(); err != nil { // 读取配置信息 读取配置信息失败
		fmt.Printf("viper Read Config failed, err:%v\n", err)
		logger.Error(map[string]interface{}{"viper Read Config failed err: ": err}, "viper Read Config Error")
		return
	}

	/* if err := viperConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("viper Read config error, config file not found")
			return
		} else {
			fmt.Println("viper Read config file error")
			return
		}
	} */

	if err := viperConfig.Unmarshal(&AppConf); err != nil {
		fmt.Printf("viper Unmarshal failed, err:%v\n", err)
		return
	}
	//fmt.Println("AppConf:", AppConf)
	logger.Info(map[string]interface{}{"AppConf Info": AppConf}, "viper Init Config Ok!!!")

	viperConfig.WatchConfig() // 对配置文件进行监视，若有改变就重新反序列到Conf中
	viperConfig.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viperConfig.Unmarshal(&AppConf); err != nil {
			fmt.Printf("viperConfig.Unmarshal failed, err:%v\n", err)
		}
	})
}
