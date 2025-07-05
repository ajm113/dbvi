package main

import "context"

type CommandHandler func(context.Context, *Editor)

type HotkeyCommand struct {
	Name        string
	Description string
	EditorModes []EditorMode
	Keys        []string
	Handler     CommandHandler // What the command does
}

var HotkeyCommandRegistry = map[string]*HotkeyCommand{}

func newHotkeyCommand(name string, description string, editorModes []EditorMode, keys []string, handler CommandHandler) *HotkeyCommand {
	return &HotkeyCommand{
		Name:        name,
		Description: description,
		EditorModes: editorModes,
		Keys:        keys,
		Handler:     handler,
	}
}

func registerHotkeyCommand(command *HotkeyCommand) {
	for _, key := range command.Keys {
		HotkeyCommandRegistry[key] = command
	}
}

type Command struct {
	Name        string
	Description string
	Command     string
	Handler     CommandHandler // What the command does
}

var CommandRegistry = map[string]*Command{}

func newCommand(name string, description string, command string, handler CommandHandler) *Command {
	return &Command{
		Name:        name,
		Description: description,
		Command:     command,
		Handler:     handler,
	}
}

func registerCommand(command *Command) {
	CommandRegistry[command.Command] = command
}
