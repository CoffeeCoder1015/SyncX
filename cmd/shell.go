package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Shell struct {
	stdin        bufio.Reader
	prompt       string
	command_list map[string]func(args *[]string)
}

// AddFunction() Adds a callable command to the shell
func (s *Shell) AddFunction(name string, action func(args *[]string)) {
	s.command_list[name] = action
}

func (s *Shell) Start() {
	s.AddFunction("help", s.Help)
	for {
		fmt.Print(s.prompt)
		text, _ := s.stdin.ReadString('\n')
		text = strings.TrimSpace(text)
		raw_cmd := strings.Split(text, " ")
		args := raw_cmd[1:]
		action, exists := s.command_list[raw_cmd[0]]
		if !exists {
			fmt.Println("Action does not exist please try again or use `help`")
			continue
		}
		action(&args)
	}
}

/*
Help() is a built in method that implements the task to display all callable commands
which the shell can execute
*/
func (s *Shell) Help(args *[]string) {
	for k := range s.command_list {
		fmt.Println(k)
	}
}

func newShell(prompt string) Shell {
	return Shell{stdin: *bufio.NewReader(os.Stdin), prompt: prompt, command_list: make(map[string]func(*[]string))}
}
