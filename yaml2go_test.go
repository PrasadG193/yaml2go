package main

import (
	"fmt"
	"go/format"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func conversion(data []byte) string {
	var d map[string]interface{}
	err := yaml.Unmarshal(data, &d)
	if err != nil {
		log.Fatal("Failed to parse input")
	}

	y2g := Yaml2Go{}
	y2g.AppendResult("type Yaml2Go struct {\n")
	y2g.Convert(d)
	y2g.AppendResult("}")

	// Convert result into go format
	_, err = format.Source([]byte(y2g.Result))
	if err != nil {
		log.Fatal("go fmt error:", err)
	}
	return string(y2g.Result)
}

func compareResults(actual []string, expected map[string]bool) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("incomplete conversion")
	}
	for _, r := range actual {
		if ok, _ := expected[r]; !ok {
			return fmt.Errorf("incomplete conversion")
		}
	}
	return nil
}

func TestConvert(t *testing.T) {
	expected := map[string]map[string]bool{
		"example1": map[string]bool{
			"type Yaml2Go struct {":                      true,
			"Array3 []float64 `yaml:\"array3\"`":         true,
			"Kind string `yaml:\"kind\"`":                true,
			"Metadata struct {":                          true,
			"Name string `yaml:\"name\"`":                true,
			"Nullfield interface{} `yaml:\"nullfield\"`": true,
			"Nestedstruct []struct {":                    true,
			"Nested struct {":                            true,

			"UnderscoreField string `yaml:\"underscore_field\"`": true,

			"Field1 []float64 `yaml:\"field1\"`":         true,
			"Field2 []bool `yaml:\"field2\"`":            true,
			"} `yaml:\"nested\"`":                        true,
			"Nested2 []struct {":                         true,
			"Nested3 struct {":                           true,
			"Field1 []int `yaml:\"field1\"`":             true,
			"Fieldt []bool `yaml:\"fieldt\"`":            true,
			"Field3 string `yaml:\"field3\"`":            true,
			"} `yaml:\"nested3\"`":                       true,
			"} `yaml:\"nested2\"`":                       true,
			"} `yaml:\"nestedstruct\"`":                  true,
			"} `yaml:\"metadata\"`":                      true,
			"Abc []struct {":                             true,
			"Def []string `yaml:\"def\"`":                true,
			"} `yaml:\"abc\"`":                           true,
			"Array1 []string `yaml:\"array1\"`":          true,
			"Array2 []int `yaml:\"array2\"`":             true,
			"IsUnderscore bool `yaml:\"is_underscore\"`": true,

			"}": true,
		},
	}

	// Test example1
	data, err := ioutil.ReadFile("tests/example1.yaml")
	if err != nil {
		log.Fatal("Failed to access tests/example1.yaml. ", err.Error())
	}
	resp := conversion(data)

	// Compare result
	err = compareResults(strings.Split(resp, "\n"), expected["example1"])
	if err != nil {
		t.Errorf("tests/example1.yaml yaml conversion incorrect")
	}
}
