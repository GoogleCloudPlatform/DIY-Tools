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
	"net/http"
	"reflect"
	"testing"
)

var parseTests = []struct {
	sURL      string
	connParam *dataConnParam
	err       error
}{
	{"https://fake.com/bp/project/dataset/view", nil, errors.New(`UnknownDataPlatform: bigquery ("bq") and firestore ("fs") are the only support platform types at this time`)},
	{"https://fake.com/bq", nil, errors.New("BadAPIRequest  Please provide a request in the following pattern\nhttps://<<hostname>>/<<data-gcp-project-target>>/platfromid/<<platform parameter 1>>/<<platform parameter 2>>")},
	{"https://fake.com/bq/project", nil, errors.New("BadAPIRequest  Please provide a request in the following pattern\nhttps://<<hostname>>/<<data-gcp-project-target>>/platfromid/<<platform parameter 1>>/<<platform parameter 2>>")},
	{"https://fake.com/bq/project/dataset/view", &dataConnParam{platform: "bq", connectionParams: []string{"project", "dataset", "view"}}, nil},
	{"https://fake.com/fs/project/collection/document", &dataConnParam{platform: "fs", connectionParams: []string{"project", "collection", "document"}}, nil},
}

func TestParseDDURL(t *testing.T) {
	for _, item := range parseTests {
		req, _ := http.NewRequest("GET", item.sURL, nil)
		pd, err := parseDDURL(req)
		if item.connParam != nil {
			if !reflect.DeepEqual(*pd, *item.connParam) {
				t.Fatalf("\nHave connection param:\n%+v\nWant connection param:\n%+v", *pd, *item.connParam)
			}
		}
		t.Logf("%+v", pd)
		if item.err != nil {
			if err.Error() != item.err.Error() {
				t.Fatalf("\nHave error:\n%s\nWant error:\n%s", err, item.err)
			}
		}
	}
}

func TestParseDataPlatfrom(t *testing.T) {
	items := parseTests[3:5]
	for _, item := range items {
		req, _ := http.NewRequest("GET", item.sURL, nil)
		pd, _ := parseDDURL(req)
		dp, err := parseDataPlatform(context.Background(), pd)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		if pd.platform == "bq" {
			if _, ok := dp.(*bqDataPlatform); !ok {
				t.Fatalf("Have:%T Want:*bqDataPlatform", dp)
			}
		}
	}

}
