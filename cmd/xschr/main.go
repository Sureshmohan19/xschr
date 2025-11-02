// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."

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
		fmt.Println("\nFor testing, run: go test ./...")
		os.Exit(1)
	}

	// TODO: Implement command routing
	fmt.Printf("Command '%s' not yet implemented\n", os.Args[1])
}
