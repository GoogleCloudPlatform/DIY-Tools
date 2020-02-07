// Copyright 2020 Google LLC

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	https://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"net/http"
	"os"
)

func init() {
	http.HandleFunc("/", getJSONData)
}

func main() {

	// Boiler plate code from the example docs for appengine hosting.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getJSONData(w http.ResponseWriter, r *http.Request) {
	// Prep auth - this will return something to be used in the query
	authHandle(w, r)

	// Parse the URL
	conParams, err := parseDDURL(r)
	if err != nil {
		handleError(w, r, err)
		return
	}

	// parse the platform interface
	pd, err := parseDataPlatform(conParams)
	if err != nil {
		handleError(w, r, err)
		return
	}

	// get the []byte results
	d := *pd
	bts, err := d.getData(r.Context())
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(bts)

}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), 500)
}
