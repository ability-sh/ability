package commander

import "fmt"

type Command struct {
	Name        string
	action      func(cmd *Command, args []string, usage bool) bool
	subCommands map[string]*Command
}

func NewCommand(name string) *Command {
	return &Command{Name: name}
}

func (C *Command) Run(args []string, usage bool) {

	if usage {

		if C.action != nil {
			C.action(C, args, usage)
		}

		for _, s := range C.subCommands {
			s.Run(args, usage)
			fmt.Println("")
		}

	} else {

		if C.action != nil {

			if C.action(C, args, usage) {
				return
			}
		}

		if len(args) > 0 && C.subCommands != nil {
			sub := C.subCommands[args[0]]
			if sub != nil {
				sub.Run(args[1:], usage)
			}
		}
	}

}

func (C *Command) SubCommand(name string) *Command {
	sub := &Command{Name: name}
	if C.Name != "" {
		sub.Name = fmt.Sprintf("%s %s", C.Name, name)
	}
	if C.subCommands == nil {
		C.subCommands = map[string]*Command{}
	}
	C.subCommands[name] = sub
	return sub
}

func (C *Command) SetAction(action func(cmd *Command, args []string, usage bool) bool) *Command {
	C.action = action
	return C
}
