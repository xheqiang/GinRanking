package models

/**
 * 之前这个包放在了dao下 调用使用dao.DB 但是dao包下就一个db链接
 * 暂时不做dao层了 放在models下 做modelbase好了
 */

import (
	"ginRanking/config"
	"ginRanking/util/logger"

	// jinzhu 导入
	/* "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" */

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"time"
)

var (
	DB  *gorm.DB
	err error
)

func init() {

	// jinzhu/gorm 的链接方式
	//DB, err = gorm.Open("mysql", config.MYSQLDB)

	MySQLUser := config.AppConf.MySQLConfig.User
	MySQLPass := config.AppConf.MySQLConfig.Password
	MySQLHost := config.AppConf.MySQLConfig.Host
	MySQLPort := config.AppConf.MySQLConfig.Port
	MySQLDB := config.AppConf.MySQLConfig.DB
	Charset := config.AppConf.MySQLConfig.Charset
	Local := config.AppConf.MySQLConfig.Local

	MySQLDSN := MySQLUser + ":" + MySQLPass + "@tcp(" + MySQLHost + ":" + MySQLPort + ")/" + MySQLDB + "?charset=" + Charset + "&parseTime=True&loc=" + Local

	// gorm 官方导入方式
	DB, err = gorm.Open(mysql.Open(MySQLDSN), &gorm.Config{})

	if err != nil {
		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}

	if DB.Error != nil {
		logger.Error(map[string]interface{}{"database connect error": DB.Error.Error()})
	}

	// jinzhu/gorm SqlDB获取方式
	// SqlDB := Db.DB()

	SqlDB, _ := DB.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	SqlDB.SetConnMaxIdleTime(time.Duration(config.AppConf.MySQLConfig.MaxIdleTime))

	// SetMaxOpenConns 设置打开数据库连接的最大数量 最大空闲链接数
	SqlDB.SetMaxOpenConns(config.AppConf.MySQLConfig.MaxOpenConns)
	SqlDB.SetMaxIdleConns(config.AppConf.MySQLConfig.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	SqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info(map[string]interface{}{"database connection initialized": "ok"})

}
