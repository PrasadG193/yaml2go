package yaml2go

import (
	"fmt"
	"go/format"
	"gopkg.in/yaml.v2"
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

// NewStruct creates new entry in StructMap result
func (yg *Yaml2Go) NewStruct(structName string, parent string) string {
	// If struct already present with the same name
	// rename struct to ParentStructname
	if _, ok := yg.StructMap[structName]; ok {
		structName = goKeyFormat(parent) + structName
	}
	yg.AppendResult(structName, fmt.Sprintf("// %s\n", structName))
	yg.StructMap[structName] += fmt.Sprintf("type %s struct {\n", structName)
	return structName
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
func (yg *Yaml2Go) Convert(structName string, data []byte) (string, error) {
	yg.StructMap = make(map[string]string)

	// Unmarshal to map[string]interface{}
	var obj map[string]interface{}
	err := yaml.Unmarshal(data, &obj)
	if err != nil {
		return "", err
	}

	yg.NewStruct("Yaml2Go", "")
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
		return "", err
	}
	return string(goFormat), nil
}

// Structify transforms map key values to struct fields
// structName : parent struct name
// k, v       : fields in the struct
func (yg *Yaml2Go) Structify(structName, k string, v interface{}, arrayElem bool) {

	if reflect.TypeOf(v) == nil || len(k) == 0 {
		yg.AppendResult(structName, fmt.Sprintf("%s interface{} `yaml:\"%s\"`\n", goKeyFormat(k), k))
		return
	}

	switch reflect.TypeOf(v).Kind() {

	// If yaml object
	case reflect.Map:
		switch val := v.(type) {
		case map[interface{}]interface{}:
			key := goKeyFormat(k)
			newKey := key
			if !arrayElem {
				// Create new structure
				newKey = yg.NewStruct(key, structName)
				yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", key, newKey, k))
			}
			// If array of yaml objects
			for k1, v1 := range val {
				if _, ok := k1.(string); ok {
					yg.Structify(newKey, k1.(string), v1, false)
				}
			}
			if !arrayElem {
				yg.AppendResult(newKey, "}\n")
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
			// Create new structure
			newKey := yg.NewStruct(key, structName)
			yg.AppendResult(structName, fmt.Sprintf("%s []%s `yaml:\"%s\"`\n", key, newKey, k))
			for _, v1 := range val {
				yg.Structify(newKey, key, v1, true)
			}
			yg.AppendResult(newKey, "}\n")
		}

	default:
		yg.AppendResult(structName, fmt.Sprintf("%s %s `yaml:\"%s\"`\n", goKeyFormat(k), reflect.TypeOf(v).String(), k))
	}
}
