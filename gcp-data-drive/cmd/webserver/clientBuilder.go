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
	"fmt"
)

// dataPlatform defines the methods needed for consumtion by the web serving handler.
type dataPlatform interface {
	// getData returns the slice of bytes that have been marshaled from the underlying data source.
	getData(ctx context.Context) ([]byte, error)
	close()
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
