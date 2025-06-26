package main

import "github.com/gdamore/tcell"

type Editor struct {
	Lines         []string
	CursorX       int
	CursorY       int
	ScrollOffsetY int
	ScrollOffsetX int
	InsertMode    bool
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
		InsertMode: false,
		screen:     screen,
	}

	editor.StatusBar = NewStatusBar(screen, editor)

	return editor
}

func (e *Editor) HandleEventKey(ek *tcell.EventKey) {

	if !e.StatusBar.StatusMode {
		e.StatusBar.HandleEventKey(ek)
		return
	}

	// General navigation that should work on all modes.
	switch ek.Key() {
	case tcell.KeyLeft:
		e.MoveCursor(-1, 0)
	case tcell.KeyRight:
		e.MoveCursor(1, 0)
	case tcell.KeyUp:
		e.MoveCursor(0, -1)
	case tcell.KeyDown:
		e.MoveCursor(0, 1)
	}

	if e.InsertMode {
		e.handleEventKeyInsertMode(ek)
	} else {
		e.handleEventKeyViewMode(ek)
	}
}

func (e *Editor) setInsertMode(insertMode bool) {
	e.InsertMode = insertMode

	if e.InsertMode {
		e.StatusBar.Command = "-- INSERT --"
	} else {
		e.StatusBar.Command = ""
	}
}

func (e *Editor) handleEventKeyViewMode(ek *tcell.EventKey) {
	switch ek.Key() {
	case tcell.KeyRune:
		switch ek.Rune() {
		case ':':
			e.StatusBar.StatusMode = false
			e.StatusBar.Command = ":"
			e.StatusBar.CursorX++
		case '/':
			e.StatusBar.StatusMode = false
			e.StatusBar.Command = "/"
			e.StatusBar.CursorX++

		// Insert mode commands
		case 'i':
			e.setInsertMode(true)
		case 'I':
			e.setInsertMode(true)
			e.SetCursor(0, e.CursorY)
		case 'o':
			e.setInsertMode(true)
			e.Lines = append(e.Lines[:e.CursorY+1], append([]string{""}, e.Lines[e.CursorY+1:]...)...)
			e.SetCursor(0, e.CursorY+1)
		case 'O':
			e.setInsertMode(true)
			e.Lines = append(e.Lines[:e.CursorY], append([]string{""}, e.Lines[e.CursorY:]...)...)
			e.SetCursor(0, e.CursorY)
		case 'a':
			e.setInsertMode(true)
		case 'A':
			e.setInsertMode(true)
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
		e.setInsertMode(false)
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
