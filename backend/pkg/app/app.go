package browserkubeapp

import (
	"context"
	"os"

	"github.com/mattn/go-colorable"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	browserkubeutil "github.com/browserkube/browserkube/pkg/util"
)

func Run(opts ...fx.Option) {
	logger := newLogger()
	opts = append(
		opts,
		fx.WithLogger(standardLogger),
		fx.Supply(logger),
		fx.Invoke(syncLogger),
	)
	app := fx.New(opts...)
	app.Run()
}

func newLogger() *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	loggerLvlStr := browserkubeutil.FirstNonEmpty(os.Getenv("LOGGER_LEVEL"), "DEBUG")
	loggerLvl, err := zapcore.ParseLevel(loggerLvlStr)
	if err != nil {
		loggerLvl = zap.DebugLevel
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(colorable.NewColorableStdout()),
		loggerLvl,
	))
	zap.ReplaceGlobals(logger)
	return logger.Sugar()
}

func syncLogger(lc fx.Lifecycle, logger *zap.SugaredLogger) {
	lc.Append(fx.Hook{OnStop: func(ctx context.Context) error {
		_ = logger.Sync()
		return nil
	}})
}

func standardLogger(logger *zap.SugaredLogger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger.Desugar()}
}
