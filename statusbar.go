package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type StatusBar struct {
	style       tcell.Style
	insertStyle tcell.Style
	normalStyle tcell.Style

	screen tcell.Screen

	Command string
	CursorX int

	editor *Editor
}

func NewStatusBar(screen tcell.Screen, editor *Editor) *StatusBar {
	return &StatusBar{
		style:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
		insertStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen),
		normalStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue),
		screen:      screen,
		editor:      editor,
	}
}

func (s *StatusBar) HandleEventKey(ek *tcell.EventKey) {
	if s.editor.EditorMode != CommandMode {
		return
	}

	switch ek.Key() {
	case tcell.KeyEscape:
		s.editor.SetEditorMode(NormalMode)
		s.Command = ""
		s.CursorX = 0
	case tcell.KeyEnter:
		s.editor.SetEditorMode(NormalMode)
		s.Command = ""
		s.CursorX = 0
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if s.CursorX > 0 {
			s.Command = s.Command[:s.CursorX-1] + s.Command[s.CursorX:]
			s.CursorX--
			break
		}
	case tcell.KeyRune:
		s.Command = s.Command[:s.CursorX] + string(ek.Rune()) + s.Command[s.CursorX:]
		s.CursorX++
	}
}

func (s *StatusBar) Draw() {
	s.drawStatus()
	s.drawCommand()
}

func (s *StatusBar) drawCommand() {
	w, h := s.screen.Size() // Get width and height

	for x := 0; x < w; x++ {
		ch := ' '
		if x < len(s.Command) {
			ch = rune(s.Command[x])
		}

		s.screen.SetContent(x, h-1, ch, nil, s.style)
	}
}

func (s *StatusBar) drawStatus() {
	w, h := s.screen.Size() // Get width and height

	mode := "UNKOWN"
	switch s.editor.EditorMode {
	case NormalMode:
		mode = "NORMAL"
	case InsertMode:
		mode = "INSERT"
	case VisualMode:
		mode = "VI-LINE"
	case CommandMode:
		mode = "COMMAND"
	}

	mode = fmt.Sprintf("  %s  ", mode)

	status := fmt.Sprintf("%s %s %d/%d:%d", mode, "[No Name]", s.editor.CursorY+1, len(s.editor.Lines), s.editor.CursorX+1)
	for x := 0; x < w; x++ {
		ch := ' '
		if x < len(status) {
			ch = rune(status[x])
		}

		style := s.style

		if x < len(mode) && s.editor.EditorMode == InsertMode {
			style = s.insertStyle
		}
		if x < len(mode) && s.editor.EditorMode != InsertMode {
			style = s.normalStyle
		}

		s.screen.SetContent(x, h-2, ch, nil, style)
	}
}
