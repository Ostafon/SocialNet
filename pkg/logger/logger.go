package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Log *zap.SugaredLogger

func Init(serviceName string) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var baseLogger *zap.Logger
	var err error

	if env == "production" {
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{"stdout"}
		cfg.InitialFields = map[string]interface{}{
			"service": serviceName,
		}
		baseLogger, err = cfg.Build()
	} else {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.OutputPaths = []string{"stdout"}
		cfg.InitialFields = map[string]interface{}{
			"service": serviceName,
		}
		baseLogger, err = cfg.Build()
	}

	if err != nil {
		fmt.Println("❌ Failed to init logger:", err)
		os.Exit(1)
	}

	Log = baseLogger.Sugar()
	Log.Infof("🚀 Logger initialized for %s", serviceName)
}

// Sync — завершает работу логгера корректно
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
