package yaml2go

import (
	"fmt"
	"go/format"
	"gopkg.in/yaml.v2"
	"log"
	"reflect"
	"strings"
)

// New creates Yaml2Go object
func New() Yaml2Go {
	return Yaml2Go{}
}

// Yaml2Go to store converted result
type Yaml2Go struct {
	StructMap map[string]string
}

// AppendResult add lines to the result
func (yg *Yaml2Go) AppendResult(structName string, line string) {
	yg.StructMap[structName] += line
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
func (yg *Yaml2Go) Convert(structName string, data []byte) string {
	yg.StructMap = make(map[string]string)

	// Unmarshal to map[string]interface{}
	var obj map[string]interface{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		log.Fatal("Failed to parse input")
	}

	yg.AppendResult("Yaml2Go", "// Yaml2Go\n")
	yg.AppendResult("Yaml2Go", "type Yaml2Go struct {\n")
	for k, v := range obj {
		yg.Structify(structName, k, v, false)
	}
	yg.AppendResult("Yaml2Go", "}\n")

	var result string
	for _, value := range yg.StructMap {
		result += fmt.Sprintf("%s\n", value)
	}
	// Convert result into go format
	goFormat, err := format.Source([]byte(result))
	if err != nil {
		log.Fatal("go fmt error:", err)
	}
	return string(goFormat)
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
