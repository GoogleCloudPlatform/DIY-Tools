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
	"strings"

	"cloud.google.com/go/firestore"
)

var errFSParams = errors.New("The url path must be in the form https://host/fs/project/collection/doc/collection/doc")

// fsDataPlatform implements getData()
type fsDataPlatform struct {

	// Pointer to the firestore client
	client *firestore.Client

	// The assembled path to the document or collection in firestore
	itemPath string

	// Indicates if the item path prepresents a firestore document or collection
	isDoc bool
}

// getData is the implementation specific to firestore for extracting
// a document of collection of documents. It satisfies the interface needed
// in the web handler
func (f *fsDataPlatform) getData(ctx context.Context) ([]byte, error) {

	// If the path is to a document, fulfill with a single doc request
	if f.isDoc {
		doc, err := f.client.Doc(f.itemPath).Get(ctx)
		if err != nil {
			return nil, err
		}
		docItem := doc.Data()
		return json.Marshal(&docItem)
	}

	// Otherwise the request is for a collection
	q := f.client.Collection(f.itemPath)

	// Get all the documents in a single read. Only a single read cost is charged.
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Create a map to hold our firestore result set.
	res := []map[string]interface{}{}

	for _, doc := range docs {
		// Adding the doc id to the result of ease of use.
		d := doc.Data()
		d["docid"] = doc.Ref.ID

		// Append the doc to the map so it can be marshaled.
		res = append(res, d)
	}
	return json.Marshal(res)
}

func getFSInterface(p *dataConnParam) (*fsDataPlatform, error) {

	// Create a platform variable
	var fsResults fsDataPlatform

	// Validate the connection params
	err := validateFSConnectionParams(p)
	if err != nil {
		return nil, err
	}

	// Setup the connection to BQ
	fsResults.client, err = firestore.NewClient(context.Background(), p.connectionParams[0])
	if err != nil {
		return nil, err
	}

	// Firestore document pattern is collection/doc/collection/doc...
	// if the item path is even then we know it is not a doc and the
	// collection get pattern applies
	fsResults.isDoc = len(p.connectionParams[1:])%2 == 0

	fsResults.itemPath = strings.Join(p.connectionParams[1:], "/")

	// Return the type
	return &fsResults, nil
}

// validateConnectionParams is a basic len check of the parameters
// TODO: Do some basic parsing to validate parameters.
func validateFSConnectionParams(p *dataConnParam) error {
	if len(p.connectionParams) < 1 {
		return errBQParams
	}
	return nil
}
