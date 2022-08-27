package commander

import (
	"fmt"
	"strings"
)

const (
	FLAG_TYPE_STRING = 0
	FLAG_TYPE_BOOL   = 1
)

type Flag struct {
	Name  string
	Usage string
	Value string
	Type  int
}

type Command struct {
	Name        string
	action      func(cmd *Command) bool
	subCommands map[string]*Command
	flagSet     map[string]*Flag
	valueSet    map[string]string
}

func NewCommand(name string) *Command {
	return &Command{Name: name}
}

func (C *Command) SetString(name string, usage string, value string) *Command {
	if C.flagSet == nil {
		C.flagSet = map[string]*Flag{}
	}
	C.flagSet[name] = &Flag{Type: FLAG_TYPE_STRING, Usage: usage, Name: name, Value: value}
	return C
}

func (C *Command) SetBool(name string, usage string) *Command {
	if C.flagSet == nil {
		C.flagSet = map[string]*Flag{}
	}
	C.flagSet[name] = &Flag{Type: FLAG_TYPE_BOOL, Usage: usage, Name: name}
	return C
}

func (C *Command) String(name string) string {

	if C.flagSet != nil {
		fs, ok := C.flagSet[name]
		if ok {
			if C.valueSet != nil {
				v, ok := C.valueSet[name]
				if ok {
					return v
				}
			}
			return fs.Value
		}
	}

	return ""
}

func (C *Command) Bool(name string) bool {

	if C.flagSet != nil {
		_, ok := C.flagSet[name]
		if ok {
			if C.valueSet != nil {
				_, ok = C.valueSet[name]
				if ok {
					return true
				}
			}
		}
	}

	return false
}

func (C *Command) Parse(args []string) {

	n := len(args)

	if C.flagSet != nil {

		C.valueSet = map[string]string{}

		for i := 0; i < n; i++ {
			key := args[i]
			if strings.HasPrefix(key, "-") {
				key = key[1:]
				v := C.flagSet[key]
				if v != nil {
					if v.Type == FLAG_TYPE_BOOL {
						C.valueSet[key] = ""
						continue
					} else if i+1 < n {
						C.valueSet[key] = args[i+1]
						i = i + 1
						continue
					}
				}
			}
		}
	}

	if C.action != nil {
		if C.action(C) {
			return
		}
	}

	if C.subCommands != nil && n > 0 {
		cmd := C.subCommands[args[0]]
		if cmd != nil {
			cmd.Parse(args[1:])
		}
	}
}

func (C *Command) Usage() {
	fmt.Printf("Usage of %s:\n", C.Name)
	for _, v := range C.flagSet {
		if v.Type == FLAG_TYPE_BOOL {
			fmt.Printf("  -%s %s\n", v.Name, v.Usage)
		} else {
			fmt.Printf("  -%s <string> %s %s\n", v.Name, v.Value, v.Usage)
		}
	}
	for _, cmd := range C.subCommands {
		cmd.Usage()
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

func (C *Command) SetAction(action func(cmd *Command) bool) *Command {
	C.action = action
	return C
}
