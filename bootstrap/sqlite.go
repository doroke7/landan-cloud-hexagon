package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"time"

	clog "github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	oSqlDB, err := oSqliteConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.SQLITE.MAX_IDLE_CONNECTIONS)

	clog.Info("[INFO] SQLITE 連線完成. ", CONFIG.SQLITE.PATH)

	return oSqliteConnectionPool, nil
}
