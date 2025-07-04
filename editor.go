package main

import (
	"unicode"

	"github.com/gdamore/tcell"
)

type EditorMode int

const (
	NormalMode EditorMode = iota
	InsertMode
	VisualMode
	CommandMode
	ExecuteMode
)

type Editor struct {
	Lines         []string
	CursorX       int
	CursorY       int
	ScrollOffsetY int
	ScrollOffsetX int
	EditorMode    EditorMode
	Width         int
	Height        int
	StatusBar     *StatusBar

	screen tcell.Screen
}

func NewEditor(screen tcell.Screen) *Editor {

	editor := &Editor{
		Lines:      []string{""},
		CursorX:    0,
		CursorY:    0,
		EditorMode: NormalMode,
		screen:     screen,
	}

	editor.StatusBar = NewStatusBar(screen, editor)

	return editor
}

func (e *Editor) HandleEventKey(ek *tcell.EventKey) {
	if e.EditorMode == CommandMode {
		e.StatusBar.HandleEventKey(ek)
		return
	}

	moveByWord := ek.Modifiers()&tcell.ModCtrl != 0

	// General navigation that should work on all modes.
	switch ek.Key() {
	case tcell.KeyLeft:
		if moveByWord {
			x := moveToPrevWord(e.Lines[e.CursorY], e.CursorX)

			if x != e.CursorX {
				e.SetCursor(x, e.CursorY)
			} else if e.CursorY > 0 {
				e.SetCursor(len(e.Lines[e.CursorY-1]), e.CursorY-1)
			}
		} else {
			e.MoveCursor(-1, 0)
		}

	case tcell.KeyRight:
		if moveByWord {
			x := moveToNextWord(e.Lines[e.CursorY], e.CursorX)

			if x != e.CursorX {
				e.SetCursor(x, e.CursorY)
			} else if e.CursorY+1 < len(e.Lines) {
				e.SetCursor(len(e.Lines[e.CursorY+1]), e.CursorY+1)
			}
		} else {
			e.MoveCursor(1, 0)
		}
	case tcell.KeyUp:
		e.MoveCursor(0, -1)
	case tcell.KeyDown:
		e.MoveCursor(0, 1)
	}

	if e.EditorMode == InsertMode {
		e.handleEventKeyInsertMode(ek)
	} else {
		e.handleEventKeyNormalMode(ek)
	}
}

func (e *Editor) SetEditorMode(editorMode EditorMode) {
	e.EditorMode = editorMode

	switch e.EditorMode {
	case InsertMode:
		e.StatusBar.Command = "-- INSERT --"
	case VisualMode:
		e.StatusBar.Command = "-- VISUAL --"
	case ExecuteMode:
		e.StatusBar.Command = "-- EXECUTE --"
	default:
		e.StatusBar.Command = ""
	}
}

func (e *Editor) handleEventKeyNormalMode(ek *tcell.EventKey) {
	switch ek.Key() {
	case tcell.KeyRune:
		switch ek.Rune() {
		case ':':
			e.SetEditorMode(CommandMode)
			e.StatusBar.Command = ":"
			e.StatusBar.CursorX++
		case '/':
			e.SetEditorMode(CommandMode)
			e.StatusBar.Command = "/"
			e.StatusBar.CursorX++

		// Insert mode commands
		case 'i':
			e.SetEditorMode(InsertMode)
		case 'I':
			e.SetEditorMode(InsertMode)
			e.SetCursor(0, e.CursorY)
		case 'o':
			e.SetEditorMode(InsertMode)
			e.Lines = append(e.Lines[:e.CursorY+1], append([]string{""}, e.Lines[e.CursorY+1:]...)...)
			e.SetCursor(0, e.CursorY+1)
		case 'O':
			e.SetEditorMode(InsertMode)
			e.Lines = append(e.Lines[:e.CursorY], append([]string{""}, e.Lines[e.CursorY:]...)...)
			e.SetCursor(0, e.CursorY)
		case 'a':
			e.SetEditorMode(InsertMode)
		case 'A':
			e.SetEditorMode(InsertMode)
			e.SetCursor(len(e.Lines[e.CursorY]), e.CursorY)

		// navigation commands
		case '0':
			e.SetCursor(0, e.CursorY)
		case '$':
			e.SetCursor(len(e.Lines[e.CursorY])-1, e.CursorY)
		case 'G':
			e.SetCursor(0, len(e.Lines)-1)
		}
	}
}

func (e *Editor) handleEventKeyInsertMode(ek *tcell.EventKey) {
	switch ek.Key() {
	case tcell.KeyEscape:
		e.SetEditorMode(NormalMode)
	case tcell.KeyEnter:
		rest := e.Lines[e.CursorY][e.CursorX:]
		e.Lines[e.CursorY] = e.Lines[e.CursorY][:e.CursorX]
		e.Lines = append(e.Lines[:e.CursorY+1], append([]string{rest}, e.Lines[e.CursorY+1:]...)...)
		e.SetCursor(0, e.CursorY+1)
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if e.CursorX > 0 {
			line := e.Lines[e.CursorY]
			e.Lines[e.CursorY] = line[:e.CursorX-1] + line[e.CursorX:]
			e.MoveCursor(-1, 0)
			break
		}

		// If we hit the end. Splice the line we are on and move it to the line above.
		if e.CursorX == 0 && e.CursorY > 0 {
			line := e.Lines[e.CursorY]

			newCursorY := e.CursorY - 1
			newCursorX := len(e.Lines[newCursorY])
			e.Lines[newCursorY] = e.Lines[newCursorY][:newCursorX] + line

			// delete the dead line and move the cursor.
			e.Lines = append(e.Lines[:e.CursorY], e.Lines[e.CursorY+1:]...)
			e.SetCursor(newCursorX, newCursorY)
		}

	case tcell.KeyRune:
		line := e.Lines[e.CursorY]
		e.Lines[e.CursorY] = line[:e.CursorX] + string(ek.Rune()) + line[e.CursorX:]
		e.MoveCursor(1, 0)
	}
}

func (e *Editor) MoveCursor(x, y int) {
	e.SetCursor(e.CursorX+x, e.CursorY+y)
}

func (e *Editor) SetCursor(x, y int) {
	if x < 0 {
		x = 0
	}

	if y < 0 {
		y = 0
	}

	e.CursorX = x
	e.CursorY = y

	if e.CursorY > len(e.Lines)-1 {
		e.CursorY = len(e.Lines) - 1
	}

	if e.CursorX > len(e.Lines[e.CursorY])-1 {
		e.CursorX = len(e.Lines[e.CursorY])
	}

	if e.CursorY < e.ScrollOffsetY {
		e.ScrollOffsetY = e.CursorY
	}

	if e.CursorY >= e.ScrollOffsetY+e.Height {
		e.ScrollOffsetY = e.CursorY - e.Height + 1
	}
}

func (e *Editor) Draw() {
	screenWidth, screenHeight := e.screen.Size()

	// We make sure the editor is aware of it's visibility
	e.Height = screenHeight - 2 // leave space for status bar
	e.Width = screenWidth

	for y := 0; y < e.Height; y++ {
		lineIndex := e.ScrollOffsetY + y
		if lineIndex >= len(e.Lines) {
			break
		}

		line := e.Lines[lineIndex]
		for x, ch := range line {
			e.screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
		}
	}

	e.StatusBar.Draw()
}

func isWordChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func moveToNextWord(s string, pos int) int {
	runes := []rune(s)
	n := len(runes)

	for pos < n && isWordChar(runes[pos]) {
		pos++
	}

	for pos < n && !isWordChar(runes[pos]) {
		pos++
	}

	return pos
}

func moveToPrevWord(s string, pos int) int {
	runes := []rune(s)

	for pos > 0 && isWordChar(runes[pos-1]) {
		pos--
	}

	for pos > 0 && !isWordChar(runes[pos-1]) {
		pos--
	}

	return pos
}
