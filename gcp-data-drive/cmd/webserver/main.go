// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
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

func main() {

	// Register the initial HTTP handler.
	http.HandleFunc("/", getJSONData)

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

	// Parse the request URL.
	conParams, err := parseDDURL(r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Parse the platform interface from the URL path.
	pd, err := parseDataPlatform(r.Context(), conParams)
	defer pd.close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Get the []byte results from the requested data platfrom.
	bts, err := pd.getData(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Setting the default content-type header to JSON.
	w.Header().Add("Content-Type", "application/json")

	// Writing the bytes to the IO writer.
	w.Write(bts)

}
