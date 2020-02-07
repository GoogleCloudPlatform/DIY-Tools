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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

var errBQParams = errors.New("The url path must be in the form https://host/bq/project/dataset/view")

// bqDataPlatform contains the nessessary information to connect ag get data from Google
// bigquery. getData() is implmented to satisify the the dataplatform interface.
type bqDataPlatform struct {
	client    *bigquery.Client
	dataQuery string
	query     *bigquery.Query
}

// getData contains the implementation detail for retriving and marshaling data from bq into JSON
// It is consumed by web server handler
func (b *bqDataPlatform) getData(ctx context.Context) ([]byte, error) {

	// Call to read to get the BQ interator
	data, err := b.query.Read(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map to hold our BQ result set.
	res := []map[string]bigquery.Value{}

	// Add the rows to a map string interface for marshaling
	for {
		row := make(map[string]bigquery.Value)
		err := data.Next(&row)

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

// getBQInterface creates and populates the bq platform client requirement and returns
// a usable getData() method
func getBQInterface(p *dataConnParam) (*bqDataPlatform, error) {

	// Create a platform variable
	var bqResults bqDataPlatform

	// Validate the connection params
	err := validateConnectionParams(p)
	if err != nil {
		return nil, err
	}

	// Setup the connection to BQ
	bqResults.client, err = bigquery.NewClient(context.Background(), p.connectionParams[0])
	if err != nil {
		return nil, err
	}

	// Create the Big Query from the parameters
	bqResults.query = bqResults.getBQQuery(p)

	// Return the type
	return &bqResults, nil
}

// getBQQuery prepares the SQL statement for the client to run.
func (b *bqDataPlatform) getBQQuery(p *dataConnParam) *bigquery.Query {
	qs := fmt.Sprintf("select * from `%s`", strings.Join(p.connectionParams, "."))
	q := b.client.Query(qs)
	q.UseStandardSQL = true
	return q
}

// validateConnectionParams is a basic len check of the parameters
// TODO: Do some basic parsing to validate parameters.
func validateConnectionParams(p *dataConnParam) error {
	if len(p.connectionParams) != 3 {
		return errBQParams
	}
	return nil
}
