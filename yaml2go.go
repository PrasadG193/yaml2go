package main

import (
	"bufio"
	"fmt"
	"go/format"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"strings"
)

const helpMsg = `yaml2go converts YAML specs to Go type definitions

Usage:
    yaml2go < /path/to/yamlspec.yaml

Examples:
    yaml2go < test/example1.yaml
    yaml2go < test/example1.yaml > example1.go
`

// Yaml2Go to store converted result
type Yaml2Go struct {
	//Result    string
	StructMap map[string]string
}

// AppendResult add lines to the result
func (yg *Yaml2Go) AppendResult(structName string, line string) {
	yg.StructMap[structName] += line
}

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

	// Unmarshal to map[string]interface{}
	var t map[string]interface{}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatal("Failed to parse input")
	}

	y2g := Yaml2Go{}
	y2g.StructMap = make(map[string]string)
	y2g.Convert("Yaml2Go", t)

	var result string
	for _, value := range y2g.StructMap {
		result += fmt.Sprintf("%s\n", value)
	}
	// Convert result into go format
	goFormat, err := format.Source([]byte(result))
	if err != nil {
		log.Fatal("go fmt error:", err)
	}
	fmt.Printf(string(goFormat))
	return
}

// removeUnderscores and camelize string
func goKeyFormat(key string) string {
	var st string
	strList := strings.Split(key, "_")
	for _, str := range strList {
		st += strings.Title(str)
	}
	if len(st) == 0 {
		st = key
	}
	return st
}

// Convert transforms map[string]interface{} to go struct
func (yg *Yaml2Go) Convert(structName string, obj map[string]interface{}) {
	yg.AppendResult("Yaml2Go", "// Yaml2Go\n")
	yg.AppendResult("Yaml2Go", "type Yaml2Go struct {\n")
	for k, v := range obj {
		yg.Structify(structName, k, v, false)
	}
	yg.AppendResult("Yaml2Go", "}\n")
}

// Structify transforms map key values to struct fields
func (yg *Yaml2Go) Structify(structName, k string, v interface{}, arrayElem bool) {
	if reflect.TypeOf(v) == nil {
		yg.AppendResult(structName, fmt.Sprintf("%s interface{} `yaml:\"%s\"`\n", goKeyFormat(k), k))
		return
	}

	switch reflect.TypeOf(v).Kind() {

	// If yaml object
	case reflect.Map:
		switch val := v.(type) {
		case map[interface{}]interface{}:
			key := goKeyFormat(k)
			if !arrayElem {
				yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", key, key, k))
				// Create new structure
				yg.AppendResult(key, fmt.Sprintf("// %s\n", key))
				yg.AppendResult(key, fmt.Sprintf("type %s struct {\n", key))
			}
			// If array of yaml objects
			for k1, v1 := range val {
				if _, ok := k1.(string); ok {
					yg.Structify(key, k1.(string), v1, false)
				}
			}
			if !arrayElem {
				yg.AppendResult(key, "}\n")
			}
		}

	// If array
	case reflect.Slice:
		val := v.([]interface{})
		if len(val) == 0 {
			return
		}
		switch val[0].(type) {

		case string, int, bool, float64:
			yg.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(val[0]), k))

		// if nested object
		case map[interface{}]interface{}:
			key := goKeyFormat(k)
			yg.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", key, key, k))
			// Create new structure
			yg.AppendResult(key, fmt.Sprintf("// %s\n", key))
			yg.AppendResult(key, fmt.Sprintf("type %s struct {\n", key))
			for _, v1 := range val {
				yg.Structify(key, key, v1, true)
			}
			yg.AppendResult(key, "}\n")
		}

	default:
		yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(v).String(), k))
	}
}
