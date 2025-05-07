package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func syncRun(config_file Config, file string) {
	data, read_file_err := os.ReadFile(config_file.Input_file)
	if read_file_err != nil {
		fmt.Println("read_file_err:", read_file_err)
		return
	}
	input_data := strings.Split(string(data), config_file.Batch_separator)
	for i, t := range input_data {
		(input_data)[i] = strings.TrimSpace(t)
	}
	//build
	raw_build_str := config_file.Build_command
	raw_build_str = strings.Replace(raw_build_str, "<file>", file, -1)

	build_cmd_str := strings.Split(raw_build_str, " ")
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
	raw_run_str := config_file.Run_command
	raw_run_str = strings.Replace(raw_run_str, "<file>", file, -1)

	run_cmd_str := strings.Split(raw_run_str, " ")
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
}
