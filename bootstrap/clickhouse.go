package bootstrap

import (
	"fmt"
	"log"
	"os"
	"time"

	clog "github.com/charmbracelet/log"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewClickhouse() (*gorm.DB, error) {

	sDSN := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%s/%s?dial_timeout=%s&read_timeout=%s",
		CONFIG.CLICKHOUSE.USER,
		CONFIG.CLICKHOUSE.PASSWORD,
		CONFIG.CLICKHOUSE.HOST,
		CONFIG.CLICKHOUSE.PORT,
		CONFIG.CLICKHOUSE.NAME,
		CONFIG.CLICKHOUSE.DIAL_TIMEOUT,
		CONFIG.CLICKHOUSE.READ_TIMEOUT,
	)

	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oClickhouseConnectionPool, err := gorm.Open(clickhouse.Open(sDSN), &gorm.Config{
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

	oSqlDB, err := oClickhouseConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.CLICKHOUSE.MAX_IDLE_CONNECTIONS)

	clog.Info("[INFO] CLICKHOUSE 連線完成. ", CONFIG.CLICKHOUSE.HOST+":"+CONFIG.CLICKHOUSE.PORT)

	return oClickhouseConnectionPool, nil
}
