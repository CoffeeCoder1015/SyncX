package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Build_command   string
	Run_command     string
	Batch_separator string
	Input_file      string
}

func loadConfig(configs *Config) {
	data, read_filer_err := os.ReadFile("./config.json")
	if read_filer_err != nil {
		error_string := pwettyPwint(read_filer_err.Error(), textProperties{Bold: true, Color: "#fc0303"})
		fmt.Println(error_string)
		return
	}
	json_unmarshal_error := json.Unmarshal(data, configs)
	if json_unmarshal_error != nil {
		fmt.Println("json_unmarshal_error:", json_unmarshal_error)
	}
}

func parseArgs(config_file Config) {
	args := os.Args
	fmt.Println(args)
	if len(args) < 2 {
		println("No arguments provided")
		return
	}
	argsWNoCommand := args[2:]

	template :=
		`{
			"build_command":"",
			"run_command":"go run <file>",
			"batch_separator":"!new",
			"input_file":"input.txt"
		}`

	command := args[1]
	switch command {
	case "init":
		// do init
		if _, err := os.Stat("config.json"); errors.Is(err, os.ErrNotExist) {
			file, _ := os.Create("config.json")
			file.WriteString(template)
		}
	case "run":
		placeHolderInBuild := strings.Contains(config_file.Build_command, "<file>")
		placeHolderInRun := strings.Contains(config_file.Run_command, "<file>")
		placeHolderFound := placeHolderInBuild || placeHolderInRun
		file := ""
		if placeHolderFound {
			if len(argsWNoCommand) == 0 {
				fmt.Println(pwettyPwint("specify a file to be compiled", textProperties{Bold: true, Color: "#fc0303"}))
				return
			}
			file = (argsWNoCommand)[0]
		}
		syncRun(config_file, file)
	}

}

func main() {

	config_file := Config{}
	loadConfig(&config_file)

	parseArgs(config_file)

	shell_name := pwettyPwint("syncx", textProperties{Color: []int{34, 105, 212}})
	shell := newShell(shell_name + " >>>")

	shell.AddFunction("cd", func(args *[]string) {
		if len(*args) == 0 {
			current_working_dir, _ := os.Getwd()
			fmt.Println(current_working_dir)
			return
		}
		os.Chdir((*args)[0])
	})

	shell.AddFunction("ls", func(args *[]string) {
		entries, _ := os.ReadDir(".")
		fmt.Println(entries)
	})

	shell.AddFunction("config", func(args *[]string) {
		loadConfig(&config_file) //calling running config command while running automatically reloads the configurations
		fmt.Println(config_file)
	})

	shell.AddFunction("run", func(args *[]string) {
		loadConfig(&config_file)
		placeHolderInBuild := strings.Contains(config_file.Build_command, "<file>")
		placeHolderInRun := strings.Contains(config_file.Run_command, "<file>")
		placeHolderFound := placeHolderInBuild || placeHolderInRun
		file := ""
		if placeHolderFound {
			if len(*args) == 0 {
				fmt.Println(pwettyPwint("specify a file to be compiled", textProperties{Bold: true, Color: "#fc0303"}))
				return
			}
			file = (*args)[0]
		}
		syncRun(config_file, file)
	})

	shell.Start()
}
