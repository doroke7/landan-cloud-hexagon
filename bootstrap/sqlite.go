package bootstrap

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewSqlite() (*gorm.DB, error) {

	if sDir := filepath.Dir(CONFIG.SQLITE.PATH); sDir != "." {
		if err := os.MkdirAll(sDir, 0o755); err != nil {
			return nil, err
		}
	}

	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oSqliteConnectionPool, err := gorm.Open(sqlite.Open(CONFIG.SQLITE.PATH), &gorm.Config{
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
			TablePrefix: CONFIG.SQLITE.PREFIX,
		},
		TranslateError: true, // 把 driver 專屬錯誤（如 UNIQUE constraint failed）轉成 gorm.ErrDuplicatedKey 等
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oSqliteConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.SQLITE.MAX_IDLE_CONNECTIONS)

	log.Info("[INFO] SQLITE 連線完成.", "addr", CONFIG.SQLITE.PATH)

	return oSqliteConnectionPool, nil
}
