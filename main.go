package main

import (
	"fmt"
	"os"
	"tropos/cmd"
)

func main() {
	fmt.Println("Starting Tropos, PID:", os.Getpid())
	cmd.Execute()
}
