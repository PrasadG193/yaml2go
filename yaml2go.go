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

// Yaml2Go to store converted result
type Yaml2Go struct {
	Result string
}

// AppendResult add lines to the result
func (yg *Yaml2Go) AppendResult(line string) {
	yg.Result += line
}

func main() {
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
	y2g.AppendResult("type Yaml2Go struct {\n")
	y2g.Convert(t)
	y2g.AppendResult("}")

	// Convert result into go format
	goFormat, err := format.Source([]byte(y2g.Result))
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

// Convert transforms map[string]interfaceP{} to go struct
func (yg *Yaml2Go) Convert(obj map[string]interface{}) {
	for k, v := range obj {
		yg.Structify(k, v, false)
	}
}

// Structify transforms map key values to struct fields
func (yg *Yaml2Go) Structify(k string, v interface{}, arrayElem bool) {
	if reflect.TypeOf(v) == nil {
		yg.AppendResult(fmt.Sprintf("%s interface{} `yaml:\"%s\"`\n", goKeyFormat(k), k))
		return
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Map:
		switch val := v.(type) {
		case map[interface{}]interface{}:
			if !arrayElem {
				yg.AppendResult(fmt.Sprintf("%s struct {\n", goKeyFormat(k)))
			}
			for k1, v1 := range val {
				if _, ok := k1.(string); ok {
					yg.Structify(k1.(string), v1, false)
				}
			}
			if !arrayElem {
				yg.AppendResult(fmt.Sprintf("} `yaml:\"%s\"`\n", k))
			}
		}

	case reflect.Slice:
		val := v.([]interface{})
		if len(val) == 0 {
			return
		}
		switch val[0].(type) {

		case string, int, bool, float64:
			yg.AppendResult(fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(val[0]), k))

		case map[interface{}]interface{}:
			yg.AppendResult(fmt.Sprintf("%s []struct {\n", goKeyFormat(k)))
			for _, v1 := range val {
				yg.Structify(goKeyFormat(k), v1, true)
			}
			yg.AppendResult(fmt.Sprintf("} `yaml:\"%s\"`\n", k))
		}

	default:
		yg.AppendResult(fmt.Sprintf("%s %s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(v).String(), k))
	}
}
