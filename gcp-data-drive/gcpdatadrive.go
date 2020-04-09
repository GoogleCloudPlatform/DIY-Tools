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
package gcpdatadrive

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// dataPlatform defines the methods needed for consumtion by the web serving handler.
type dataPlatform interface {
	// getData returns the slice of bytes that have been marshaled from the underlying data source.
	getData(ctx context.Context) ([]byte, error)
	close() error
}

func GetJSONData(w http.ResponseWriter, r *http.Request) {
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

// parseDataPlatform detects the requested data platform and returns an interface to specified data platform.
func parseDataPlatform(ctx context.Context, p *dataConnParam) (dataPlatform, error) {
	switch p.platform {
	case "bq":
		return newBQPlatform(ctx, p)

	case "fs":
		return newFSPlatform(ctx, p)

	}

	return nil, fmt.Errorf(`unknown data platform %q: bigquery ("bq") and firestore ("fs") supported`, p.platform)
}

// dataConnParam provides parsed parameters from the requested URL path.
type dataConnParam struct {
	// platfrom is a charter indicator of the target platform. Accepted values are (bq,fs).
	platform string

	// connectionParams is the remaining path from the url request split on a "/" charter.
	connectionParams []string
}

// parseDDURL detects and shapes the data platfrom request.
func parseDDURL(r *http.Request) (*dataConnParam, error) {
	// Trimming the leading prefix and split the path in to an array.
	location := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(location) < 3 {
		return nil, errors.New("BadAPIRequest  Please provide a request in the following pattern\nhttps://<<hostname>>/<<data-gcp-project-target>>/platfromid/<<platform parameter 1>>/<<platform parameter 2>>")
	}

	// This switch statement is used to allow easy implementation of additional dat platform providers.
	switch location[0] {
	case "bq", "fs":
		return &dataConnParam{
			platform:         location[0],
			connectionParams: location[1:],
		}, nil

	}
	return nil, errors.New(`UnknownDataPlatform: bigquery ("bq") and firestore ("fs") are the only support platform types at this time`)
}
