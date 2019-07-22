# yaml2go
Converts YAML specs into Go type definitions


[![asciicast](https://asciinema.org/a/bCu9sI4j77r2Ut8U8pnTXZoU1.svg)](https://asciinema.org/a/bCu9sI4j77r2Ut8U8pnTXZoU1)


## Usage
```bash
$ yaml2go < example.yaml
```
e.g

```bash
$ cat example1.yaml
kind: test
metadata:
  name: cluster
  nullfield:
  nestedstruct:
  - nested:
      underscore_field: value
      field1:
      - 44.5
      - 43.6
      field2:
      - true
      - false
    nested2:
      - nested3:
          field1:
          - 44
          - 43
          fieldt:
          - true
          - false
          field3: value
abc:
  - def:
    - black
    - white
array1:
  - "string1"
  - "string2"
array2:
  - 2
  - 6
array3:
  - 3.14
  - 5.12
is_underscore: true
```

```bash
$ $GOPATH/bin/yaml2go < example1.yaml
type Yaml2Go struct {
        Kind     string `yaml:"kind"`
        Metadata struct {
                Name         string      `yaml:"name"`
                Nullfield    interface{} `yaml:"nullfield"`
                Nestedstruct []struct {
                        Nested struct {
                                UnderscoreField string    `yaml:"underscore_field"`
                                Field1          []float64 `yaml:"field1"`
                                Field2          []bool    `yaml:"field2"`
                        } `yaml:"nested"`
                        Nested2 []struct {
                                Nested3 struct {
                                        Field1 []int  `yaml:"field1"`
                                        Fieldt []bool `yaml:"fieldt"`
                                        Field3 string `yaml:"field3"`
                                } `yaml:"nested3"`
                        } `yaml:"nested2"`
                } `yaml:"nestedstruct"`
        } `yaml:"metadata"`
        Abc []struct {
                Def []string `yaml:"def"`
        } `yaml:"abc"`
        Array1       []string  `yaml:"array1"`
        Array2       []int     `yaml:"array2"`
        Array3       []float64 `yaml:"array3"`
        IsUnderscore bool      `yaml:"is_underscore"`
}
```

## Installation

#### Step 1: Clone the repo
```bash
$ git clone https://github.com/PrasadG193/yaml2go.git $GOPATH/src/github.com/PrasadG193/yaml2go
```

#### Step 2: Build binary using go install
```bash
$ go install $GOPATH/src/github.com/PrasadG193/yaml2go
```

#### Step 3: Convert YAML specs into Go type

```bash
$ cat example1.yaml
kind: test
metadata:
  name: cluster
  nullfield:
  nestedstruct:
  - nested:
      underscore_field: value
      field1:
      - 44.5
      - 43.6
      field2:
      - true
      - false
    nested2:
      - nested3:
          field1:
          - 44
          - 43
          fieldt:
          - true
          - false
          field3: value
abc:
  - def:
    - black
    - white
array1:
  - "string1"
  - "string2"
array2:
  - 2
  - 6
array3:
  - 3.14
  - 5.12
is_underscore: true
```

```bash
$ $GOPATH/bin/yaml2go < example1.yaml
type Yaml2Go struct {
        Kind     string `yaml:"kind"`
        Metadata struct {
                Name         string      `yaml:"name"`
                Nullfield    interface{} `yaml:"nullfield"`
                Nestedstruct []struct {
                        Nested struct {
                                UnderscoreField string    `yaml:"underscore_field"`
                                Field1          []float64 `yaml:"field1"`
                                Field2          []bool    `yaml:"field2"`
                        } `yaml:"nested"`
                        Nested2 []struct {
                                Nested3 struct {
                                        Field1 []int  `yaml:"field1"`
                                        Fieldt []bool `yaml:"fieldt"`
                                        Field3 string `yaml:"field3"`
                                } `yaml:"nested3"`
                        } `yaml:"nested2"`
                } `yaml:"nestedstruct"`
        } `yaml:"metadata"`
        Abc []struct {
                Def []string `yaml:"def"`
        } `yaml:"abc"`
        Array1       []string  `yaml:"array1"`
        Array2       []int     `yaml:"array2"`
        Array3       []float64 `yaml:"array3"`
        IsUnderscore bool      `yaml:"is_underscore"`
}
```
