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

	StatusMode bool
	Command    string
	CursorX    int

	editor *Editor
}

func NewStatusBar(screen tcell.Screen, editor *Editor) *StatusBar {
	return &StatusBar{
		style:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGray),
		insertStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen),
		normalStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue),
		screen:      screen,
		editor:      editor,
		StatusMode:  true,
	}
}

func (s *StatusBar) HandleEventKey(ek *tcell.EventKey) {
	if s.StatusMode {
		return
	}

	switch ek.Key() {
	case tcell.KeyEscape:
		s.StatusMode = true
		s.Command = ""
		s.CursorX = 0
	case tcell.KeyEnter:
		s.StatusMode = true
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
	if s.StatusMode {
		s.drawStatus()
		return
	}

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

	mode := "  NORMAL  "
	if s.editor.InsertMode {
		mode = "  INSERT  "
	}

	status := fmt.Sprintf("%s %s %d/%d:%d", mode, "[No Name]", s.editor.CursorY+1, len(s.editor.Lines), s.editor.CursorX+1)
	for x := 0; x < w; x++ {
		ch := ' '
		if x < len(status) {
			ch = rune(status[x])
		}

		style := s.style

		if s.StatusMode {
			if x < len(mode) && s.editor.InsertMode {
				style = s.insertStyle
			}
			if x < len(mode) && !s.editor.InsertMode {
				style = s.normalStyle
			}
		}

		s.screen.SetContent(x, h-1, ch, nil, style)
	}
}
