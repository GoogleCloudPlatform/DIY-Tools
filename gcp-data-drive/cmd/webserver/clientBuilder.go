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
	"context"
)

// dataPlatform defines the methods needed for consumtion by the web serving handler
type dataPlatform interface {
	// getData returns the array of bytes that have been marshaled from the underlying data source
	getData(ctx context.Context) ([]byte, error)
}

// parseDataPlatform detects the requested data platform and returns an interface
// pointer to the web handler.
func parseDataPlatform(p *dataConnParam) (*dataPlatform, error) {
	var plat dataPlatform
	var err error

	switch p.platform {
	case "bq":
		plat, err = getBQInterface(p)
		return &plat, err
	case "fs":
		plat, err = getFSInterface(p)
		return &plat, err
	}

	return nil, errUnknowDataPlatform
}
