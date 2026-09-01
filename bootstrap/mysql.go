package bootstrap

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DDL:
// CREATE TABLE users (
//   id   INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   name VARCHAR(255) NOT NULL
// );

func NewMysql() (*gorm.DB, error) {

	var sHost, sPort string
	if len(CONFIG.MYSQL.WRITE.HOSTS) > 0 {
		sHost = CONFIG.MYSQL.WRITE.HOSTS[0]
	}
	if len(CONFIG.MYSQL.WRITE.PORTS) > 0 {
		sPort = CONFIG.MYSQL.WRITE.PORTS[0]
	}

	sDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		CONFIG.MYSQL.USERNAME,
		CONFIG.MYSQL.PASSWORD,
		sHost,
		sPort,
		CONFIG.MYSQL.NAME,
		CONFIG.MYSQL.CHARSET,
	)
	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oMysqlConnectionPool, err := gorm.Open(mysql.Open(sDSN), &gorm.Config{
		Logger: logger.New(
			gormWriter{}, // 輸出到標準輸出（見 gorm_writer.go）
			logger.Config{
				SlowThreshold:             time.Second, // 慢查詢閾值
				LogLevel:                  oLogLevel,   // 日誌級別：Silent, Error, Warn, Info
				IgnoreRecordNotFoundError: true,        // 是否忽略 ErrRecordNotFound 錯誤
				ParameterizedQueries:      false,       // 是否在日誌中顯示參數值（設為 false 會顯示具體數值）
				Colorful:                  true,        // 是否啟用彩色字體
			},
		),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: CONFIG.MYSQL.PREFIX,
		},
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oMysqlConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.MYSQL.MAX_IDLE_CONNECTIONS)

	log.Info("[INFO] MYSQL 連線完成.", "addr", sHost+":"+sPort)

	return oMysqlConnectionPool, nil
}
