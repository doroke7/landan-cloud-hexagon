package bootstrap

import (
	"fmt"
	"time"

	clog "github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GaussDB 源自 PostgreSQL，走同一套 wire protocol，所以直接沿用 gorm.io/driver/postgres，
// 差異只在連線資訊（CONFIG.GAUSSDB）跟目標叢集是 GaussDB 而不是 PostgreSQL。
func NewGaussdb() (*gorm.DB, error) {

	sDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		CONFIG.GAUSSDB.HOST,
		CONFIG.GAUSSDB.PORT,
		CONFIG.GAUSSDB.USER,
		CONFIG.GAUSSDB.PASSWORD,
		CONFIG.GAUSSDB.NAME,
		CONFIG.GAUSSDB.SSLMODE,
	)

	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oGaussdbConnectionPool, err := gorm.Open(postgres.Open(sDSN), &gorm.Config{
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
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oGaussdbConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.GAUSSDB.MAX_IDLE_CONNECTIONS)

	clog.Info("[INFO] GAUSSDB 連線完成.", "addr", CONFIG.GAUSSDB.HOST+":"+CONFIG.GAUSSDB.PORT)

	return oGaussdbConnectionPool, nil
}
