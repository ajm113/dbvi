package main

import "fmt"

type cliArguments struct {
	Files    []string
	Flags    map[rune]bool
	Commands []cliCommand
}

type cliCommand struct {
	Name  string
	Value string
}

func parseFlags(args []string) (*cliArguments, error) {
	args = args[1:]

	cliArgs := &cliArguments{}

	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "--" {
			cliArgs.Files = append(cliArgs.Files, args[i+1:]...)
			break
		}

		if len(arg) > 1 && arg[0] == '-' && arg[1] != '-' {
			for _, c := range arg[1:] {
				switch c {
				// TODO: Add flags here.
				case 'h':
					cliArgs.Commands = append(cliArgs.Commands, cliCommand{Name: "help"})
				default:
					return nil, fmt.Errorf("unknown option argument: \"%c\"", c)
				}
			}
		} else if len(arg) > 1 && arg[0] == '-' && arg[1] == '-' {
			switch arg[2:] {
			case "help":
				cliArgs.Commands = append(cliArgs.Commands, cliCommand{Name: "help"})
			default:
				return nil, fmt.Errorf("unknown option argument: \"%s\"", arg)
			}
		} else {
			cliArgs.Files = append(cliArgs.Files, arg)
		}
		i++
	}

	return cliArgs, nil
}
