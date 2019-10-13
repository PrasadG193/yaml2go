package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PrasadG193/yaml2go"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const API_VERSION = "v1"
const PORT = "8080"

func main() {
	router := httprouter.New()
	log.Printf("server started accepting requests on port=%s..\n", PORT)
	router.POST(fmt.Sprintf("/%s/convert", API_VERSION), handleConvert)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleConvert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Bad Request. Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	// Create yaml2go object and invoke Convert()
	y2g := yaml2go.New()
	result, err := y2g.Convert("Yaml2Go", []byte(data))
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Bad Request. Error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	io.WriteString(w, result)
}
