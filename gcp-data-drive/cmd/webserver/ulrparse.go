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
	"net/http"
	"strings"
)

// dataConnParam is used for parsed parameter consumption by the requested platform
type dataConnParam struct {
	// platfrom is a 2 charter indicator of the target platform. Accepted values are (bq,fs)
	platform string

	// connectionParams is the remaining path from the url request split on a "/"
	// charter
	connectionParams []string
}

// parseDDURL detects and shapes the data platfrom request.
func parseDDURL(r *http.Request) (*dataConnParam, error) {
	// Trimming the leading prefix and splitting the path in to an array
	location := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(location) < 1 {
		return nil, errBadURLPattern
	}

	// This switch statement provide and open implementation for adding new platform providers.
	switch location[0] {
	case "bq", "fs":
		return &dataConnParam{
			platform:         location[0],
			connectionParams: location[1:],
		}, nil

	}
	return nil, errUnknowDataPlatform
}
