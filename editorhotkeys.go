package main

import "context"

func setDefaultHotkeys(e *Editor) {
	// Insert hotkeys
	registerHotkeyCommand(newHotkeyCommand(
		"Insert",
		"Enters insert mode at current cursor position",
		[]EditorMode{NormalMode},
		[]string{"i"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Insert Before Non-Blank Line",
		"Enters insert mode at the first instance of a non-blank line of the current cursor's Y position",
		[]EditorMode{NormalMode},
		[]string{"I"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
			// TODO move me to the first non-blank character.
			e.SetCursor(0, e.CursorY)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Insert New Line",
		"Enters insert mode with a new line +1 of Y cursor",
		[]EditorMode{NormalMode},
		[]string{"o"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
			e.Lines = append(e.Lines[:e.CursorY], append([]string{""}, e.Lines[e.CursorY:]...)...)
			e.SetCursor(0, e.CursorY)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Insert New Line Current Position",
		"Enters insert mode with a new line at current cursor",
		[]EditorMode{NormalMode},
		[]string{"o"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
			e.Lines = append(e.Lines[:e.CursorY], append([]string{""}, e.Lines[e.CursorY:]...)...)
			e.SetCursor(0, e.CursorY)
		},
	))

	registerHotkeyCommand(newHotkeyCommand(
		"Insert after cursor",
		"Enters insert mode at +1 of cursor's x position",
		[]EditorMode{NormalMode},
		[]string{"a"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
			e.MoveCursor(1, 0)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Insert at end line",
		"Enters insert mode at end of the current line",
		[]EditorMode{NormalMode},
		[]string{"A"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(InsertMode)
			e.SetCursor(len(e.Lines[e.CursorY]), e.CursorY)
		},
	))

	// Entering visual mode
	registerHotkeyCommand(newHotkeyCommand(
		"Visual",
		"Enters visual mode",
		[]EditorMode{NormalMode, VisualMode},
		[]string{"v"},
		func(_ context.Context, e *Editor) {
			if e.EditorMode == NormalMode {
				e.SetEditorMode(VisualMode)
				e.CursorStartX = e.CursorX
				e.CursorStartY = e.CursorY
			} else {
				e.SetEditorMode(NormalMode)
			}
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Visual Line",
		"Enters visual line mode",
		[]EditorMode{NormalMode},
		[]string{"V"},
		func(_ context.Context, e *Editor) {
			e.SetEditorMode(VisualLineMode)
			e.CursorStartX = 0
			e.CursorStartY = e.CursorY
			e.SetCursor(len(e.Lines[e.CursorY]), e.CursorY)
		},
	))

	// navigation
	registerHotkeyCommand(newHotkeyCommand(
		"Move To First Character",
		"Moves cursor to the first character of a given line",
		[]EditorMode{NormalMode},
		[]string{"0"},
		func(_ context.Context, e *Editor) {
			e.SetCursor(0, e.CursorY)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Move To Last Character",
		"Moves cursor to the last character of a given line",
		[]EditorMode{NormalMode},
		[]string{"$"},
		func(_ context.Context, e *Editor) {
			e.SetCursor(len(e.Lines[e.CursorY])-1, e.CursorY)
		},
	))
	registerHotkeyCommand(newHotkeyCommand(
		"Move To Last Character of File",
		"Moves cursor to the last line of a file",
		[]EditorMode{NormalMode},
		[]string{"G"},
		func(_ context.Context, e *Editor) {
			e.SetCursor(0, len(e.Lines)-1)
		},
	))
}
