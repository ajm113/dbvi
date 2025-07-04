package main

import (
	"errors"
	"os"

	"github.com/ajm113/dbvi/config"
	"github.com/gdamore/tcell"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	screen tcell.Screen
	log    *zap.SugaredLogger
	editor *Editor
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init() {
	// TODO: Move me to a tmp dir?
	logger, err := setupLogger("dbvi.log")
	if err != nil {
		println("failed creating dbvi.log", err)
		os.Exit(1)
	}

	defer logger.Sync() // flushes buffer, if any
	a.log = logger.Sugar()

	configPath, err := config.FindDefault()

	// TODO: Have this not fail completely?
	if errors.Is(err, config.ErrConfigNotFound) {
		a.log.Fatal("config not found")
	}

	if err != nil {
		a.log.Fatal("unexpected error finding config", zap.Any("error", err))
	}

	a.log.Debugf("loading config: %s", configPath)
	_, err = config.Load(configPath)
	if err != nil {
		a.log.Fatal("unexpected error loading config", zap.Any("error", err))
	}

	a.log.Info("loaded config")

	a.screen, err = tcell.NewScreen()
	if err != nil {
		a.log.Fatal("unexpected error creating tcell screen", zap.Any("error", err))
	}
	err = a.screen.Init()
	if err != nil {
		a.log.Fatal("unexpected error init tcell screen", zap.Any("error", err))
	}

	a.editor = NewEditor(a.screen)
}

func (a *App) Run() error {
	defer func() {
		a.screen.Fini()
	}()

	a.draw()

	for {
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				a.screen.Fini()
				os.Exit(0)
			}
			a.editor.HandleEventKey(ev)
		case *tcell.EventResize:
			a.screen.Sync()
		}
		a.draw()
	}
}

func (a *App) draw() {
	a.screen.Clear()
	a.editor.Draw()

	if a.editor.EditorMode == CommandMode {
		a.screen.ShowCursor(a.editor.StatusBar.CursorX, a.editor.Height+1)
	} else {
		a.screen.ShowCursor(a.editor.CursorX, a.editor.CursorY-a.editor.ScrollOffsetY)
	}

	a.screen.Show()
}

func main() {
	app := NewApp()
	app.Init()
	err := app.Run()
	if err != nil {
		os.Exit(1)
	}
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
