package bootstrap

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Oracle 沒有 MySQL/PostgreSQL 那樣的 wire protocol 相容性，
// 所以用 github.com/godoes/gorm-oracle（底層是純 Go 實作的 go-ora，不需要 OCI client）。
func NewOracle() (*gorm.DB, error) {

	iPort, err := strconv.Atoi(CONFIG.ORACLE.PORT)
	if err != nil {
		return nil, err
	}

	sDSN := oracle.BuildUrl(
		CONFIG.ORACLE.HOST,
		iPort,
		CONFIG.ORACLE.SERVICE,
		CONFIG.ORACLE.USER,
		CONFIG.ORACLE.PASSWORD,
		nil,
	)

	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oOracleConnectionPool, err := gorm.Open(oracle.Open(sDSN), &gorm.Config{
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

	oSqlDB, err := oOracleConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.ORACLE.MAX_IDLE_CONNECTIONS)

	return oOracleConnectionPool, nil
}
