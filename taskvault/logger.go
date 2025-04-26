package taskvault

import (
	"io"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapOnce sync.Once

func InitLogger(logLevel string, node string) *zap.SugaredLogger {
	var zapLogger *zap.Logger
	var err error

	level := zapcore.InfoLevel
	if parsedLevel, err := zapcore.ParseLevel(logLevel); err == nil {
		level = parsedLevel
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: level == zapcore.DebugLevel,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.CapitalColorLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err = cfg.Build()
	if err != nil {
		panic("failed to build zap logger: " + err.Error())
	}

	sugar := zapLogger.Sugar().With("node", node)

	zapOnce.Do(func() {
		if level == zapcore.DebugLevel {
			gin.DefaultWriter = os.Stdout
			gin.SetMode(gin.DebugMode)
		} else {
			gin.DefaultWriter = io.Discard
			gin.SetMode(gin.ReleaseMode)
		}
	})

	return sugar
}
