package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
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

func main() {
	config_file := Config{}
	loadConfig(&config_file)

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
		data, read_file_err := os.ReadFile(config_file.Input_file)
		if read_file_err != nil {
			fmt.Println("read_file_err:", read_file_err)
			return
		}
		input_data := strings.Split(string(data), config_file.Batch_separator)
		for i, t := range input_data {
			(input_data)[i] = strings.TrimSpace(t)
		}
		file := (*args)[0]
		//build
		build_cmd_str := strings.Split(strings.Replace(config_file.Build_command, "<file>", file, -1), " ")
		Build := exec.Command(build_cmd_str[0], build_cmd_str[1:]...)
		compiler_err_txt, _ := Build.CombinedOutput()
		if len(compiler_err_txt) > 0 {
			fmt.Println(pwettyPwint("'run' canceled as compilation issue occurred", textProperties{Bold: true, Color: "#fc0303"}))
			fmt.Println(string(compiler_err_txt))
			return
		}
		//run
		var wg sync.WaitGroup
		mut := &sync.Mutex{}
		run_cmd_str := strings.Split(strings.Replace(config_file.Run_command, "<file>", file, -1), " ")
		results := make([]string, len(input_data))
		for i, t := range input_data {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				var output bytes.Buffer
				Run := exec.Command(run_cmd_str[0], run_cmd_str[1:]...)
				Run.Stdout = &output
				Run.Stderr = &output
				Run.Stdin = strings.NewReader(t)
				runerr := Run.Run()
				if runerr != nil {
					fmt.Println("runerr:", runerr)
				}
				mut.Lock()
				results[id] = output.String()
				mut.Unlock()
			}(i)
		}
		wg.Wait()
		for i, result := range results {
			fmt.Printf(pwettyPwint("Case %d\n", textProperties{Color: "#f542d7", Bold: true, Underline: true}), i)
			fmt.Print(result)
		}
	})

	shell.Start()
}
