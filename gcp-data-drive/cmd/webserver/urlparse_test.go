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
	"reflect"
	"testing"
)

func TestUrlNegative(t *testing.T) {
	badUrls := []string{
		"https://fake.com/",
		"fake.com",
	}

	for _, u := range badUrls {
		// Create a request test from the URL
		req, err := http.NewRequest("GET", u, nil)

		_, err = parseDDURL(req)

		// We are expecting an error to pass the test.
		if err == nil {
			t.Fatalf("Expecting Error")
		}
	}
}

func TestUrlProviders(t *testing.T) {
	urls := []string{"http://host/bq/project/dataset/view", "https://host/bq/project/dataset/view"}

	respRes := dataConnParam{
		platform:         "bq",
		connectionParams: []string{"project", "dataset", "view"},
	}

	for _, u := range urls {
		// Create a request test from the URL
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}
		con, err := parseDDURL(req)
		// We are expecting an error to pass the test.
		if err != nil {
			t.Fatalf("%v - %s", con, err.Error())
		}

		if !reflect.DeepEqual(*con, respRes) {
			t.Fatalf("Have :%+v Want:%+v", *con, respRes)
		}

	}
}
