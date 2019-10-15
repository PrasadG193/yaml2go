package yaml2go

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleConvert(w http.ResponseWriter, r *http.Request) {
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
	y2g := New()
	result, err := y2g.Convert("Yaml2Go", []byte(data))
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("Bad Request. Error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	io.WriteString(w, result)
}
