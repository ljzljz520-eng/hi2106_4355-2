package main

import (
	"flag"
	"fmt"
	"os"

	"traininggrades/internal/cli"
	"traininggrades/internal/grades"
)

func main() {
	dataPath := flag.String("data", "grades.json", "path used by save and load operations")
	fixturePath := flag.String("fixture", "", "optional JSON file loaded before the menu starts")
	flag.Parse()

	service := grades.NewService(grades.NewMemoryStore())
	if *fixturePath != "" {
		if err := service.Load(*fixturePath); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	runner := cli.NewRunner(os.Stdin, os.Stdout, service, *dataPath)
	if err := runner.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
