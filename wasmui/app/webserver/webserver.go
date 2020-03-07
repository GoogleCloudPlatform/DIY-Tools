package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	http.HandleFunc("/css/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request for %s on css", r.URL.Path)
		bts, err := ioutil.ReadFile(fmt.Sprintf("site/css/%s", path.Base(r.URL.Path)))
		if err != nil {
			log.Println(string(bts))
		}
		w.Header().Add("content-type", "text/css")
		w.Write(bts)
		return
	})

	http.HandleFunc("/js/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request for %s on javascript", r.URL.Path)

		bts, err := ioutil.ReadFile(fmt.Sprintf("site/js/%s", path.Base(r.URL.Path)))
		if err != nil {
			log.Println(string(bts))
		}
		w.Header().Add("content-type", "text/javascript")
		w.Write(bts)
		return
	})

	http.HandleFunc("/wasm/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request for %s on wasm", r.URL.Path)
		bts, err := ioutil.ReadFile("site/wasm/app.wasm")
		if err != nil {
			log.Println(string(bts))
			http.Error(w, "Internal Server", http.StatusInternalServerError)
		}
		w.Header().Add("content-type", "application/wasm")
		w.Write(bts)
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request for %s on catchall", r.URL.Path)
		bts, err := ioutil.ReadFile("site/wasm_exec.html")
		if err != nil {
			log.Println(string(bts))
			http.Error(w, "Internal Server", http.StatusInternalServerError)
		}
		w.Header().Add("content-type", "text/html")
		w.Write(bts)
		return
	})

	log.Println("Listening...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
