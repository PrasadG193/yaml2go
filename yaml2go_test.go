package yaml2go

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func conversion(data []byte) (string, error) {
	y2g := New()
	result, err := y2g.Convert("Yaml2Go", data)

	return result, err
}

func compareResults(actual []string, expected map[string]bool) error {
	for _, r := range actual {
		if len(r) == 0 {
			continue
		}
		if ok, _ := expected[r]; !ok {
			fmt.Printf("Couldn't find: %s %d\n", r, len(r))
			return fmt.Errorf("incomplete conversion")
		}
	}
	return nil
}

func TestConvert(t *testing.T) {
	expected := map[string]map[string]bool{
		"example1": map[string]bool{
			"// Yaml2Go":            true,
			"type Yaml2Go struct {": true,
			"	Kind         string    `yaml:\"kind\"`": true,
			"	Metadata     Metadata  `yaml:\"metadata\"`": true,
			"	Abc          []Abc     `yaml:\"abc\"`": true,
			"	Array1       []string  `yaml:\"array1\"`": true,
			"	Array2       []int     `yaml:\"array2\"`": true,
			"	Array3       []float64 `yaml:\"array3\"`": true,
			"	IsUnderscore bool      `yaml:\"is_underscore\"`": true,
			"// Metadata":            true,
			"type Metadata struct {": true,
			"	Name         string         `yaml:\"name\"`": true,
			"	Nullfield    interface{}    `yaml:\"nullfield\"`": true,
			"	Nestedstruct []Nestedstruct `yaml:\"nestedstruct\"`": true,
			"// Nested3":            true,
			"type Nested3 struct {": true,
			"	Field1 []int  `yaml:\"field1\"`": true,
			"	Fieldt []bool `yaml:\"fieldt\"`": true,
			"	Field3 string `yaml:\"field3\"`": true,
			"// Abc":            true,
			"type Abc struct {": true,
			"	Def []string `yaml:\"def\"`": true,
			"// Nestedstruct":            true,
			"type Nestedstruct struct {": true,
			"	Nested  Nested    `yaml:\"nested\"`": true,
			"	Nested2 []Nested2 `yaml:\"nested2\"`": true,
			"// Nested":            true,
			"type Nested struct {": true,
			"	UnderscoreField string    `yaml:\"underscore_field\"`": true,
			"	Field1          []float64 `yaml:\"field1\"`": true,
			"	Field2          []bool    `yaml:\"field2\"`": true,
			"// Nested2":            true,
			"type Nested2 struct {": true,
			"	Nested3 Nested3 `yaml:\"nested3\"`": true,
			"}": true,
		},
	}

	// Test example1
	data, err := ioutil.ReadFile("testdata/example1.yaml")
	if err != nil {
		log.Fatal("Failed to access tests/example1.yaml. ", err.Error())
	}
	resp, err := conversion(data)
	assert.Nil(t, err)

	// Compare result
	err = compareResults(strings.Split(resp, "\n"), expected["example1"])
	if err != nil {
		t.Errorf("tests/example1.yaml yaml conversion incorrect")
	}
}
