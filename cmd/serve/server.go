package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PrasadG193/yaml2go"
)

const API_VERSION = "v1"
const PORT = "8080"

func main() {
	log.Printf("server started accepting requests on port=%s..\n", PORT)
	http.HandleFunc(fmt.Sprintf("/%s/convert", API_VERSION), yaml2go.HandleConvert)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
