package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// 初始化具備分割功能的 Logger
	logger := InitLogger()
	defer logger.Sync()

	r := gin.New()

	// 模擬 middleware 記錄日誌
	r.Use(func(c *gin.Context) {
		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.Run(":8080")
}

func InitLogger() *zap.Logger {
	// 配置 Lumberjack 進行日誌分割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日誌路徑
		MaxSize:    100,              // 單個文件最大 100MB，超過就分割
		MaxBackups: 10,               // 保留最近 10 個舊檔案
		MaxAge:     30,               // 檔案保留 30 天
		Compress:   true,             // 舊檔案自動 Gzip 壓縮，節省硬碟空間
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		// 同時輸出到檔案（帶分割功能）與控制台
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	return zap.New(core)
}
