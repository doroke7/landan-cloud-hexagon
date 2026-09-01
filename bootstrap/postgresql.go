package bootstrap

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewPostgresql() (*gorm.DB, error) {

	sDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		CONFIG.POSTGRESQL.HOST,
		CONFIG.POSTGRESQL.PORT,
		CONFIG.POSTGRESQL.USERNAME,
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
			TablePrefix: CONFIG.POSTGRESQL.PREFIX,
		},
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oPostgresqlConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.POSTGRESQL.MAX_IDLE_CONNECTIONS)

	log.Info("[INFO] POSTGRESQL 連線完成.", "addr", CONFIG.POSTGRESQL.HOST+":"+CONFIG.POSTGRESQL.PORT)

	return oPostgresqlConnectionPool, nil
}
