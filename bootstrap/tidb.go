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

// TiDB 跟 MySQL 走同一套 wire protocol，所以直接沿用 gorm.io/driver/mysql，
// 差異只在連線資訊（CONFIG.TIDB）跟目標叢集是 TiDB 而不是 MySQL。
func NewTidb() (*gorm.DB, error) {

	var sHost, sPort string
	if len(CONFIG.TIDB.HOSTS) > 0 {
		sHost = CONFIG.TIDB.HOSTS[0]
	}
	if len(CONFIG.TIDB.PORTS) > 0 {
		sPort = CONFIG.TIDB.PORTS[0]
	}

	sDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		CONFIG.TIDB.USERNAME,
		CONFIG.TIDB.PASSWORD,
		sHost,
		sPort,
		CONFIG.TIDB.NAME,
		CONFIG.TIDB.CHARSET,
	)
	var oLogLevel logger.LogLevel = logger.Error
	if CONFIG.DEFAULT.DEBUG {
		oLogLevel = logger.Info
	}

	oTidbConnectionPool, err := gorm.Open(mysql.Open(sDSN), &gorm.Config{
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
			TablePrefix: CONFIG.TIDB.PREFIX, // 例如所有表都加上 sys_ 前綴
		},
	})
	if err != nil {
		return nil, err
	}

	oSqlDB, err := oTidbConnectionPool.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.TIDB.MAX_IDLE_CONNECTIONS)

	log.Info("[INFO] TIDB 連線完成.", "addr", sHost+":"+sPort)

	return oTidbConnectionPool, nil
}
