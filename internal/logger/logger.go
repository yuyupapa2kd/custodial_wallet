package logger

import (
	"custodial-vault/configs"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger
}

func Init() *Logger {
	// logger := zap.Must(zap.NewProduction())
	// logger := createLogger(true, )
	logger := createLogger()

	defer logger.Sync()

	logger.Info("zap logger initiated!!!")
	return &Logger{
		zapLogger: logger,
	}
}

func (lg *Logger) LevelError(err error, c *gin.Context) {
	lg.zapLogger.Error(err.Error(),
		zap.Int("status_code", c.Writer.Status()),
		zap.String("req_method", c.Request.Method),
		zap.String("req_uri", c.Request.RequestURI),
		zap.String("client_ip", c.ClientIP()),
	)
}

func (lg *Logger) LevelInfo(c *gin.Context) {
	lg.zapLogger.Info("succeed",
		zap.Int("status_code", c.Writer.Status()),
		zap.String("req_method", c.Request.Method),
		zap.String("req_uri", c.Request.RequestURI),
		zap.String("client_ip", c.ClientIP()),
	)
}

func (lg *Logger) VaultHealthCheck(resp *resty.Response, req_uri string, initialized bool, sealed bool, serverTime int) {
	signal := "good"
	if resp.StatusCode() != 200 {
		signal = "bad"
	}
	lg.zapLogger.Info(signal,
		zap.Int("status_code", resp.StatusCode()),
		zap.String("req_method", "GET"),
		zap.String("req_uri", req_uri),
		zap.Bool("vault.initialized", initialized),
		zap.Bool("vault.sealed", sealed),
		zap.Int("vault.serverTime", serverTime),
	)
}

// func createLogger(d bool, f *os.File) *zap.Logger {
// 	encoderCfg := zap.NewProductionEncoderConfig()
// 	encoderCfg.TimeKey = "timestamp"
// 	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
// 	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
// 	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

// 	level := zap.InfoLevel
// 	if d {
// 		level = zap.DebugLevel
// 	}

// 	core := zapcore.NewTee(
// 		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
// 		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
// 	)

// 	l := zap.New(core)

// 	return l
// }

func createLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	// fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	// consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stdout",
			configs.RuntimeConf.LogPath.Info,
		},
		ErrorOutputPaths: []string{
			"stderr",
			configs.RuntimeConf.LogPath.Error,
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getegid(),
		},
	}

	return zap.Must(config.Build())
}
