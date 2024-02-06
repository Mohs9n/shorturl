package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v3"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func BoltHandler(db *bolt.DB, bname string, fallback http.Handler) (http.HandlerFunc, error) {
	m := make(map[string]string)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bname))
		c := b.Cursor()

		for k, v := c.First(); k!= nil; k, v = c.Next(){
			m[string(k)] = string(v)
		}
		return nil
	})
	fmt.Println(m)
	return MapHandler(m, fallback,), nil
}

func parseYAML(yml []byte) ([]map[string]string, error) {
	var urls []map[string]string
	err := yaml.Unmarshal(yml, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	for _, url := range urls {
		pathsToUrls[url["path"]] = url["url"]
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	urls, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	for _, url := range urls {
		pathsToUrls[url["path"]] = url["url"]
	}
	fmt.Println(pathsToUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parseJSON(jsn []byte) ([]map[string]string, error) {
	var urls []map[string]string
	err := json.Unmarshal(jsn, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}