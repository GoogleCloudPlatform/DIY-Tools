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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// bqDataPlatform contains the necessary information to connect and get data from Bigquery platfrom.
type bqDataPlatform struct {
	// client is a pointer to a BQ client.
	client *bigquery.Client

	// dataQuery is the base query string in ANSI SQL.
	dataQuery string

	// query is  a pointer to the bigquery query struct which is composed from the dataQuery.
	query *bigquery.Query
}

// getData contains the implementation detail for retriving and marshaling data from Bigqueyr into JSON.
func (b *bqDataPlatform) getData(ctx context.Context) ([]byte, error) {
	// Call the read function to get the BQ interator of the Bigquery rows.
	it, err := b.query.Read(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map to hold our Bigquery results.
	res := []map[string]bigquery.Value{}

	// Add the Bigquery rows to a map string interface for marshaling.
	// TODO: This implementation builds and map in memory. The dataset size must fit in memory. Consider
	// providing callback fulfillment for large datasets leverging pub/sub and GCS.
	for {
		row := make(map[string]bigquery.Value)
		err := it.Next(&row)

		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}
		res = append(res, row)
	}
	b.client.Close()
	return json.Marshal(&res)
}

// getBQInterface creates and populates the Bigquery platform client requirementa and returns
// a usable getData() method.
func newBQPlatform(p *dataConnParam) (*bqDataPlatform, error) {
	// Validate the connection params and return and error if they are not compatible.
	if err := validateConnectionParams(p); err != nil {
		return nil, err
	}

	// Create the Bigquery client.
	c, err := bigquery.NewClient(p.requestContext, p.connectionParams[0])
	if err != nil {
		return nil, err
	}

	// Create an ANSI SQL Query string from the HTTP request path.
	qs := fmt.Sprintf("select * from `%s`", strings.Join(p.connectionParams, "."))

	// Create the BQ query
	q := c.Query(qs)

	// Set the standard SQL option
	q.UseStandardSQL = true

	return &bqDataPlatform{
		query:  q,
		client: c,
	}, nil

}

// validateConnectionParams is a basic len check of the parameters
// TODO: Add additional complex parsing to check the parameters.
func validateConnectionParams(p *dataConnParam) error {
	// A basic check to make sure we have at least 3 parameters to work with.
	if len(p.connectionParams) != 3 {
		return errors.New("the url path must be in the form https://host/bq/project/dataset/view")
	}
	return nil
}
