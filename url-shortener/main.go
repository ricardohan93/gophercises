package main

import (
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"net/http"
)

func main() {
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := mapHandlerFunc(pathsToUrls, mux)

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := yAMLHandlerFunc([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func mapHandlerFunc(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the path from the incoming URL
		path := r.URL.Path
		// check if the path exist in our mapping
		if dest, ok := pathsToUrls[path]; ok {
			// if the path exists, we redirect the user to that URL (dest)
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func yAMLHandlerFunc(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// parse the yaml
	pathUrls, err := parseYaml(yamlBytes)
	if err != nil {
		return nil, err
	}
	// convert yaml array into map
	pathsToUrls := buildMap(pathUrls)
	// return a mapHandler
	return mapHandlerFunc(pathsToUrls, fallback), nil
}

func buildMap(pathUrls []pathURL) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, pu := range pathUrls {
		pathsToUrls[pu.Path] = pu.URL
	}
	return pathsToUrls
}

func parseYaml(data []byte) ([]pathURL, error) {
	var pathUrls []pathURL
	err := yaml.Unmarshal(data, &pathUrls)
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
