package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

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

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parseYaml(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func parseYaml(yml []byte) (map[string]string, error) {
	type yamlSchema struct {
		Path string `yaml:"path"`
		Url  string `yaml:"url"`
	}

	var parsedYaml []yamlSchema
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
