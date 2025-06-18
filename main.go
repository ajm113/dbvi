package main

import (
	"errors"
	"os"

	"github.com/ajm113/dbvi/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	// TODO: Move me to a tmp dir?
	logger, err := setupLogger("dbvi.log")
	if err != nil {
		println("failed creating dbvi.log", err)
		os.Exit(1)
	}

	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	configPath, err := config.FindDefault()

	// TODO: Have this not fail completely?
	if errors.Is(err, config.ErrConfigNotFound) {
		sugar.Fatal("config not found")
	}

	if err != nil {
		sugar.Fatal("unexpected error finding config", zap.Any("error", err))
	}

	sugar.Debugf("loading config: %s", configPath)
	_, err = config.Load(configPath)
	if err != nil {
		sugar.Fatal("unexpected error loading config", zap.Any("error", err))
	}

	sugar.Info("loaded config")
}

func setupLogger(logFile string) (*zap.Logger, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	// encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder // optional: colorized output
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(file),
		zap.InfoLevel,
	)
	return zap.New(core), nil
}
