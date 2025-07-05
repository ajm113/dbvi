package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ajm113/dbvi/utils"
	"github.com/gdamore/tcell"
)

type EditorMode int

const (
	NormalMode EditorMode = iota
	InsertMode
	VisualMode
	VisualLineMode
	CommandMode
	ExecuteMode
)

type Editor struct {
	Lines              []string
	Clipboard          []string
	ClipboardMultiline bool
	CursorX            int
	CursorY            int
	CursorStartX       int
	CursorStartY       int
	ScrollOffsetY      int
	ScrollOffsetX      int
	EditorMode         EditorMode
	Width              int
	Height             int
	StatusBar          *StatusBar

	screen        tcell.Screen
	normalStyle   tcell.Style
	selectedStyle tcell.Style
	executedStyle tcell.Style
	bufferedKeys  string
}

func NewEditor(screen tcell.Screen) *Editor {

	editor := &Editor{
		Lines:         []string{""},
		CursorX:       0,
		CursorY:       0,
		EditorMode:    NormalMode,
		screen:        screen,
		normalStyle:   tcell.StyleDefault,
		selectedStyle: tcell.StyleDefault.Foreground(tcell.ColorGrey).Background(tcell.ColorWhite),
	}

	editor.StatusBar = NewStatusBar(screen, editor)
	setDefaultHotkeys(editor)

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
	case ':':
		if e.EditorMode != InsertMode && e.EditorMode != CommandMode {
			e.SetEditorMode(CommandMode)
			e.StatusBar.Command = ":"
			e.StatusBar.CursorX++
		}
	case '/':
		if e.EditorMode != InsertMode && e.EditorMode != CommandMode {
			e.SetEditorMode(CommandMode)
			e.StatusBar.Command = "/"
			e.StatusBar.CursorX++
		}
	case tcell.KeyEscape:
		e.SetEditorMode(NormalMode)
	case tcell.KeyLeft:
		if moveByWord {
			x := utils.MoveToPrevWord(e.Lines[e.CursorY], e.CursorX)

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
			x := utils.MoveToNextWord(e.Lines[e.CursorY], e.CursorX)

			if x != e.CursorX {
				e.SetCursor(x, e.CursorY)
			} else if e.CursorY+1 < len(e.Lines) {
				e.SetCursor(len(e.Lines[e.CursorY+1]), e.CursorY+1)
			}
		} else {
			e.MoveCursor(1, 0)
		}
	case tcell.KeyUp:

		if e.EditorMode != VisualLineMode {
			e.MoveCursor(0, -1)
		} else {
			e.SetCursor(len(e.Lines[e.CursorY]), e.CursorY-1)
		}
	case tcell.KeyDown:
		if e.EditorMode != VisualLineMode {
			e.MoveCursor(0, 1)
		} else {
			e.SetCursor(len(e.Lines[e.CursorY]), e.CursorY+1)
		}
	}

	switch e.EditorMode {
	case InsertMode:
		e.handleEventKeyInsertMode(ek)
	}

	e.handleHotkeys(ek)
}

func (e *Editor) SetEditorMode(editorMode EditorMode) {
	e.EditorMode = editorMode
	e.bufferedKeys = ""

	switch e.EditorMode {
	case InsertMode:
		e.StatusBar.Command = "-- INSERT --"
	case VisualMode:
		e.StatusBar.Command = "-- VISUAL --"
	case VisualLineMode:
		e.StatusBar.Command = "-- VISUAL LINE --"
	case ExecuteMode:
		e.StatusBar.Command = "-- EXECUTE --"
	default:
		e.StatusBar.Command = ""
	}
}

func (e *Editor) handleHotkeys(ek *tcell.EventKey) {
	if e.EditorMode == InsertMode {
		return
	}

	e.bufferedKeys += eventKeyToString(ek)
	if cmd, ok := HotkeyCommandRegistry[e.bufferedKeys]; ok {

		if len(cmd.EditorModes) == 0 {
			cmd.Handler(context.Background(), e)
		} else {
			for _, m := range cmd.EditorModes {
				if m == e.EditorMode {
					cmd.Handler(context.Background(), e)
					break
				}
			}
		}
		e.bufferedKeys = ""
	}

	// TOOD: This is prone to bugs and should be updated to actually count # of key presses in the string...
	if len(strings.Split(e.bufferedKeys, "+")) > 1 || len(e.bufferedKeys) > 1 {
		e.bufferedKeys = ""
	}
}

func (e *Editor) handleEventKeyInsertMode(ek *tcell.EventKey) {
	switch ek.Key() {
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

			style := e.normalStyle

			if e.isSelected(x, y) {
				style = e.selectedStyle
			}

			e.screen.SetContent(x, y, ch, nil, style)
		}
	}

	e.StatusBar.Draw()
}

func (e *Editor) isSelected(x, y int) bool {
	if e.EditorMode != VisualMode && e.EditorMode != VisualLineMode {
		return false
	}

	startX := e.CursorStartX
	startY := e.CursorStartY

	endX := e.CursorX
	endY := e.CursorY

	// If for some reason these are the wrong way, flip em.
	if y < startY || (y == startY && x < startX) {
		startX, endX = endX, startX
		startY, endY = endY, startY
	}

	if y < startY || y > endY {
		return false
	}

	if y == startY && y == endY {
		return x >= startX && x <= endX
	}

	if y == startY {
		return x >= startX
	}

	if y == endY {
		return x <= endX
	}

	return true
}

func eventKeyToString(ev *tcell.EventKey) string {
	mod := ev.Modifiers()
	var prefix string

	if mod&tcell.ModCtrl != 0 {
		prefix += "Ctrl+"
	}
	if mod&tcell.ModAlt != 0 {
		prefix += "Alt+"
	}
	if mod&tcell.ModShift != 0 {
		prefix += "Shift+"
	}

	switch ev.Key() {
	case tcell.KeyRune:
		return prefix + string(ev.Rune())
	case tcell.KeyEsc:
		return prefix + "Esc"
	case tcell.KeyEnter:
		return prefix + "Enter"
	case tcell.KeyTab:
		return prefix + "Tab"
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return prefix + "Backspace"
	case tcell.KeyUp:
		return prefix + "Up"
	case tcell.KeyDown:
		return prefix + "Down"
	case tcell.KeyLeft:
		return prefix + "Left"
	case tcell.KeyRight:
		return prefix + "Right"
	case tcell.KeyHome:
		return prefix + "Home"
	case tcell.KeyEnd:
		return prefix + "End"
	case tcell.KeyDelete:
		return prefix + "Delete"
	case tcell.KeyInsert:
		return prefix + "Insert"
	case tcell.KeyPgUp:
		return prefix + "PageUp"
	case tcell.KeyPgDn:
		return prefix + "PageDown"
	case tcell.KeyCtrlSpace:
		return "Ctrl+Space"
	default:
		return fmt.Sprintf("%s[%v]", prefix, ev.Key())
	}
}
