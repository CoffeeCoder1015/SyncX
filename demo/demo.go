package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func squared() {
	reader := bufio.NewReader(os.Stdin)
	str, _ := reader.ReadString('\n')
	num, _ := strconv.Atoi(strings.TrimSpace(str))
	fmt.Println(num * num)
}

func main() {
	squared()
}
