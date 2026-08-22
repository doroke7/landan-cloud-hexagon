package bootstrap

import (
	"fmt"
	"log"
	"os"
	"time"

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
	if len(CONFIG.DATABASE.WRITE.HOSTS) > 0 {
		sHost = CONFIG.DATABASE.WRITE.HOSTS[0]
	}
	if len(CONFIG.DATABASE.WRITE.PORTS) > 0 {
		sPort = CONFIG.DATABASE.WRITE.PORTS[0]
	}

	sDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		CONFIG.DATABASE.USER,
		CONFIG.DATABASE.PASSWORD,
		sHost,
		sPort,
		CONFIG.DATABASE.NAME,
		CONFIG.DATABASE.CHARSET,
	)
	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oMysqlConnectionPool, err := gorm.Open(mysql.Open(sDSN), &gorm.Config{
		// TranslateError: true, // 把 MySQL 原生 error（例如 1062 重複 key）轉成 gorm.ErrDuplicatedKey 等標準化 error，全專案統一判斷
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // 輸出到標準輸出
			logger.Config{
				SlowThreshold:             time.Second, // 慢查詢閾值
				LogLevel:                  oLogLevel,   // 日誌級別：Silent, Error, Warn, Info
				IgnoreRecordNotFoundError: true,        // 是否忽略 ErrRecordNotFound 錯誤
				ParameterizedQueries:      false,       // 是否在日誌中顯示參數值（設為 false 會顯示具體數值）
				Colorful:                  true,        // 是否啟用彩色字體
			},
		),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "",
		},
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oMysqlConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.DATABASE.MAX_IDLE_CONNECTIONS)

	return oMysqlConnectionPool, nil
}
