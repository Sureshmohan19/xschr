package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("XSchr - Experimental Job Scheduler")
	fmt.Println("Version: 1.0.0")
	
	if len(os.Args) < 2 {
		fmt.Println("\nUsage: xschr <command>")
		fmt.Println("\nCommands:")
		fmt.Println("  controller  - Start the scheduler controller")
		fmt.Println("  worker      - Start a worker node")
		fmt.Println("  submit      - Submit a job")
		os.Exit(1)
	}
}

