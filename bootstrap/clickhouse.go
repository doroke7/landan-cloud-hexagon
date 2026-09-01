package bootstrap

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewClickhouse() (*gorm.DB, error) {

	sDSN := fmt.Sprintf(
		"clickhouse://%s:%s@%s:%s/%s?dial_timeout=%s&read_timeout=%s",
		CONFIG.CLICKHOUSE.USERNAME,
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
			TablePrefix: CONFIG.CLICKHOUSE.PREFIX,
		},
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oClickhouseConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.CLICKHOUSE.MAX_IDLE_CONNECTIONS)

	log.Info("[INFO] CLICKHOUSE 連線完成.", "addr", CONFIG.CLICKHOUSE.HOST+":"+CONFIG.CLICKHOUSE.PORT)

	return oClickhouseConnectionPool, nil
}
