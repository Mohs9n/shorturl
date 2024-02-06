package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mohs9n/shorturl/handler"
	"github.com/boltdb/bolt"
)

func main() {
	mux := defaultMux()

	yamlFile := flag.String("yaml", "", "a YAML file with path/url pairs")
	jsonFile := flag.String("json", "", "a JSON file with path:url pairs")
	boltFile := flag.String("bolt", "paths.db", "a BoltDB database file")
	flag.Parse()

	
	var yml []byte
	var jsn []byte
	var err error

	db, err := bolt.Open(*boltFile, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b :=tx.Bucket([]byte("Paths"))
		
		b.Put([]byte("/hello"), []byte("https://x.com"))
		return nil
	})

	if *jsonFile != "" {
		jsn, err = os.ReadFile(*jsonFile)
		if err != nil {
			panic(err)
		}
	} else {
		jsn = []byte(`
[]
`)
	}

	if *yamlFile == "" {
		yml = []byte(
			`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	} else {
		// // Open the YAML file
		// ymlFile, err := os.Open(*yamlFile)
		// if err != nil {
		// 	panic(err)
		// }
		// defer ymlFile.Close()

		yml, err = os.ReadFile(*yamlFile)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("%s\n",yml)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v3",
	}
	mapHandler := handler.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	// 	yaml := `
	// - path: /urlshort
	//   url: https://github.com/gophercises/urlshort
	// - path: /urlshort-final
	//   url: https://github.com/gophercises/urlshort/tree/solution
	// `
	yamlHandler, err := handler.YAMLHandler(yml, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := handler.JSONHandler(jsn, yamlHandler)
	if err != nil {
		panic(err)
	}

	boltHandler, err := handler.BoltHandler(db, "Paths", jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", boltHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}