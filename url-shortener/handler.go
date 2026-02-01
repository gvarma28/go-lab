package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type schema struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if destUrl, ok := pathsToUrls[url]; ok {
			http.Redirect(w, r, destUrl, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseJson(jsonData)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(yml []byte) (map[string]string, error) {
	var parsedYaml []schema
	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string)
	for _, v := range parsedYaml {
		urlMap[v.Path] = v.Url
	}
	return urlMap, nil
}

func parseJson(jsonData []byte) (map[string]string, error) {
	var parsedJson []schema
	err := json.Unmarshal(jsonData, &parsedJson)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]string)
	for _, v := range parsedJson {
		urlMap[v.Path] = v.Url
	}
	return urlMap, nil
}
