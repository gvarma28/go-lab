package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	urlshort "github.com/gvarma28/go-lab/url-shortener"
)

func main() {
	yamlFile := flag.String("yamlFile", "urls.yaml", "path to urls yaml file")
	jsonFile := flag.String("jsonFile", "urls.json", "path to urls json file")
	flag.Parse()

	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlData, err := os.ReadFile(*yamlFile)
	if err != nil {
		panic(err)
	}
	yamlHandler, err := urlshort.YAMLHandler([]byte(yamlData), mapHandler)
	if err != nil {
		panic(err)
	}

	jsonData, err := os.ReadFile(*jsonFile)
	if err != nil {
		panic(err)
	}
	jsonHandler, err := urlshort.YAMLHandler([]byte(jsonData), yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
