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
	"context"
	"errors"
	"net/http"
	"strings"
)

// dataConnParam provides parsed parameters from the requested URL path.
type dataConnParam struct {
	// platfrom is a charter indicator of the target platform. Accepted values are (bq,fs).
	platform string

	// connectionParams is the remaining path from the url request split on a "/" charter.
	connectionParams []string

	// requestContext is the orginal http request context.
	requestContext context.Context
}

// parseDDURL detects and shapes the data platfrom request.
func parseDDURL(r *http.Request) (*dataConnParam, error) {

	// Trimming the leading prefix and split the path in to an array.
	location := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(location) < 1 {
		return nil, errors.New("BadAPIRequest  Please provide a request in the following pattern\nhttps://<<hostname>>/<<data-gcp-project-target>>/platfromid/<<platform parameter 1>>/<<platform parameter 2>>")
	}

	// This switch statement is used to allow easy implementation of additional dat platform providers.
	switch location[0] {
	case "bq", "fs":
		return &dataConnParam{
			platform:         location[0],
			connectionParams: location[1:],
			requestContext:   r.Context(),
		}, nil

	}
	return nil, errors.New(`UnknownDataPlatform: bigquery ("bq") and firestore ("fs") are the only support platform types at this time`)
}
