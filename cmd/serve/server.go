package main

import (
	"fmt"
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
	router.POST(fmt.Sprintf("/%s/convert", API_VERSION), yaml2go.HandleConvert)

	handler := cors.Default().Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

