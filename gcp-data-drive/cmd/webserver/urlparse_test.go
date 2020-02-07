package main

import (
	"net/http"
	"reflect"
	"testing"
)

func TestUrlNegative(t *testing.T) {
	badUrls := []string{
		"http://fake.com/",
		"http://fake.com",
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
