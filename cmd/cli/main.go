package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/PrasadG193/yaml2go"
)

const helpMsg = `yaml2go converts YAML specs to Go type definitions

Usage:
    yaml2go < /path/to/yamlspec.yaml

Examples:
    yaml2go < test/example1.yaml
    yaml2go < test/example1.yaml > example1.go
`

func printHelp(f string) {
	helpArgs := []string{"-h", "--help", "help"}
	for _, m := range helpArgs {
		if f == m {
			fmt.Println(helpMsg)
			os.Exit(0)
		}
	}
}

func main() {
	// Read args
	if len(os.Args) > 1 {
		printHelp(os.Args[1])
	}

	// Read input from the console
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading input:", err)
	}

	// Create yaml2go object and invoke Convert()
	y2g := yaml2go.New()
	result, err := y2g.Convert("Yaml2Go", []byte(data))
	if err != nil {
		log.Fatal("Invalid YAML")
	}

	fmt.Printf(result)
	return
}
