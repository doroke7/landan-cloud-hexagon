package bootstrap

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresql() (*gorm.DB, error) {

	sDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		CONFIG.POSTGRESQL.HOST,
		CONFIG.POSTGRESQL.PORT,
		CONFIG.POSTGRESQL.USER,
		CONFIG.POSTGRESQL.PASSWORD,
		CONFIG.POSTGRESQL.NAME,
		CONFIG.POSTGRESQL.SSLMODE,
	)

	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oPostgresqlConnectionPool, err := gorm.Open(postgres.Open(sDSN), &gorm.Config{
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
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oPostgresqlConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.POSTGRESQL.MAX_IDLE_CONNECTIONS)

	return oPostgresqlConnectionPool, nil
}
