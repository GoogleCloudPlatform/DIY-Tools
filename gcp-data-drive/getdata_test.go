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
	"net/http"
	"reflect"
	"testing"
)

func TestParseDDURL(t *testing.T) {
	var tests = []struct {
		in    string
		out   *dataConnParam
		isErr bool
	}{
		{"https://example.com/bp/project/dataset/view",
			nil,
			true,
		},
		{"https://example.com/bq",
			nil,
			false,
		},
		{"https://example.com/bq/project",
			nil,
			true,
		},

		{"https://example.com/bq/project/dataset/view",
			&dataConnParam{platform: "bq", connectionParams: []string{"project", "dataset", "view"}},
			false,
		},
		{"https://example.com/fs/project/collection/document",
			&dataConnParam{platform: "fs", connectionParams: []string{"project", "collection", "document"}},
			false,
		},
	}

	for pos, item := range tests {
		req, err := http.NewRequest("GET", item.in, nil)
		if err != nil {
			t.Errorf("parseDDURL(%v): error creating a fake http request", item.in)
		}

		pd, err := parseDDURL(req)
		if item.out != nil {
			if !reflect.DeepEqual(*pd, *item.out) {
				t.Errorf("parseDDURL(%v)\nHave connection param:\n%+v\nWant connection param:\n%+v", item.in, *pd, *item.out)
			}
		}
		if item.isErr {
			if err == nil {
				t.Errorf("parseDDURL(%v) test:%v, An error was expected but no error was returned", item.in, pos)
			}
		}
	}
}

func TestParseDataPlatfrom(t *testing.T) {
	var platformDetectTests = []struct {
		in string
	}{
		{"https://example.com/bq/project/dataset/view"},
		{"https://example.com/fs/project/collection/document"},
	}

	for _, item := range platformDetectTests {
		req, err := http.NewRequest("GET", item.in, nil)
		if err != nil {
			t.Errorf("parseDataPlatform() Error creating fake http request.")
		}

		pd, err := parseDDURL(req)
		if err != nil {
			t.Errorf("parseDataPlatform() Error parsing the parameters in the URL")
		}

		have, err := parseDataPlatform(context.Background(), pd)
		if err != nil {
			t.Errorf("parseDataPlatform(): Expecting no errors in the test but have %v", err)
		}

		if pd.platform == "bq" {
			if _, ok := have.(*bqDataPlatform); !ok {
				t.Errorf("parseDataPlatform(context,%+v) = %T Want:*bqDataPlatform", pd, have)
			}
		}

		if pd.platform == "fs" {
			if _, ok := have.(*fsDataPlatform); !ok {
				t.Errorf("parseDataPlatform(context,%+v) = %T Want:*fsDataPlatform", pd, have)
			}
		}
	}

}
