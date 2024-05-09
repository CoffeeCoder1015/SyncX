#   SYNCX  
**SYNC**hronous E**X**ecution is a tool for testing programs that take `stdin` inputs 

It does this in parallel leading to **Blazingly fast testing cycles!!! 🔥🔥**

## Features
- Build & Run all in one go
- SYNCX runs all the test cases over a program in parallel leading to shorter waiting time

## Usage
SYNCX is a shell where there are command which can be executed to do the testing

[make sure these files have been properly configured](#before-starting-some-files-need-to-be-present-and-setup-for-syncx-to-run-correctly)

### Standard directory commands
```
cd <dir | no dir>
ls
```
calling `cd <no dir>` will print out the current working directory

### Config
```
config 
```
each call of this command will try to read [config.json](/demo/config.json) in the local directory 
and update SYNCX configs for program testing

### Run
```
run <file | no file>
```
the build / run commands are all abstracted behing a single `run`

A specific file can be specified if a placeholder `<file>` tag is present [build_command](#build_command) or [run_command](#run_command)

- if a `<file>` placeholder is not found, SYNCX assumes the file has been specified in the build command rather than in `run`

## Before starting, some files need to be present and setup for syncx to run correctly
- config.json 
- a input file
### demo [config.json](/demo/config.json)
this file needs to be in the current working directory of SYNCX

```json
{
    "build_command":"",
    "run_command":"go run <file>",
    "batch_separator":"!new",
    "input_file":"input.txt"
}
```
#### build_command
SYNCX can build the program for you before executing it

`build_command` allows you to provide a build command to compile the program into a executable.

This filed would be used for languages like c++ where there is no equivalent to `go run` or `cargo run`  

- There is a default placeholder tag, `<file>`, which can be used to specify the file to be compiled during runtime
#### run_command
`run_command` specifies the shell command to execute to start the specified program 

- There is a default placeholder tag, `<file>`, which can be used to specify the file to be ran during runtime

#### input_file
The `stdin` input which SYNCX will send to the program being ran will be sourced from a file

#### batch_separator
`batch_separator` specifies a substring the the `input_file` which will be interpreted as 
a separator for a new set of input data to be executed

### demo [input.txt](/demo/input.txt)
in the input.txt file each input case is separated by the [batch_separator](#batch_separator) `!new`
```
1            
!new
2
!new
...
```

[demo.go](/demo/demo.go) Takes in a single number and prints out the result squared

Then SYNCX would run [demo.go](/demo/demo.go) on each of the input cases in sync

#### Results (The case numbers would be colored pink)
```
Case 0
1
Case 1
4
...
```

## Build from source
```powershell
cd /cmd
go build .
```