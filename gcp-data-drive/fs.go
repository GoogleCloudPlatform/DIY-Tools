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
	"encoding/json"
	"errors"
	"strings"

	"cloud.google.com/go/firestore"
)

// fsDataPlatform contains the necessary information to connect and get data from Firestore platfrom.
type fsDataPlatform struct {

	// client is a pointer to the firestore client.
	client *firestore.Client

	// itemPath is an assembled path to the document or collection in firestore.
	itemPath string

	// isDoc indicates if the item path prepresents a firestore document or collection
	isDoc bool
}

// getData is the implementation specific to firestore for extracting a document of collection of documents.
func (f *fsDataPlatform) getData(ctx context.Context) ([]byte, error) {
	// If the path is to a document, fulfill the request with the document.
	if f.isDoc {
		doc, err := f.client.Doc(f.itemPath).Get(ctx)
		if err != nil {
			return nil, err
		}
		docItem := doc.Data()
		return json.Marshal(&docItem)
	}

	// Otherwise the request is for a collection.
	q := f.client.Collection(f.itemPath)

	// Get all the documents in a single read. Only a single read is charged.
	docs, err := q.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	// Create a slice of maps to hold the firestore result set.
	res := []map[string]interface{}{}

	for _, doc := range docs {
		// Adding the doc id to the result for ease of use.
		d := doc.Data()
		d["docid"] = doc.Ref.ID

		// Append the doc to the map so it can be marshaled.
		res = append(res, d)
	}

	return json.Marshal(res)
}

// close will close the firestore client connection
func (f *fsDataPlatform) close() error {
	if err := f.client.Close(); err != nil {
		return err
	}
	return nil
}

func newFSPlatform(ctx context.Context, p *dataConnParam) (*fsDataPlatform, error) {
	// Validate the connection parameters.
	if err := validateFSConnectionParams(p); err != nil {
		return nil, err
	}

	// Create the connection to Firestore.
	client, err := firestore.NewClient(ctx, p.connectionParams[0])

	return &fsDataPlatform{
		client: client,

		// Firestore document pattern is collection/doc/collection/doc... if the item path is even then we know
		// it is not a doc and the collection get logic applies.
		isDoc: len(p.connectionParams[1:])%2 == 0,

		// Join the Firestore doc path from the parsed parameters.
		itemPath: strings.Join(p.connectionParams[1:], "/"),
	}, err

}

// validateConnectionParams is a basic len check of the parameters.
// TODO: Do additional parsing to validate parameters.
func validateFSConnectionParams(p *dataConnParam) error {
	if len(p.connectionParams) < 1 {
		return errors.New("the url path must be in the form https://host/fs/project/collection/doc/collection/doc")
	}
	return nil
}
