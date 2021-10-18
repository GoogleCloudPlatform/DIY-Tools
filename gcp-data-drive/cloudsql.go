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
	"fmt"
	"log"

	"encoding/json"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// fsDataPlatform contains the necessary information to connect and get data from Firestore platfrom.
type sqlDataPlatform struct {

	// client is a pointer to the firestore client.
	client *sql.DB

	// itemPath is an assembled path to the document or collection in firestore.
	query string

	// isDoc indicates if the item path prepresents a firestore document or collection
	//isDoc bool
}

// getData is the implementation specific to firestore for extracting a document of collection of documents.
func (s *sqlDataPlatform) getData(ctx context.Context) ([]byte, error) {

	rows, err := s.client.Query(s.query)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %v", err)
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Count the columns
	count := len(columns)

	// Make a overall results collector. Will be marshaled later
	results := make([]map[string]interface{}, 0)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row error: %v", err)
	}

	return json.Marshal(results)
}

// close will close the firestore client connection
func (s *sqlDataPlatform) close() error {
	// if s.client == nil {
	// 	return nil
	// }

	// if err := s.client.Close(); err != nil {
	// 	return err
	// }
	return nil
}

func newSQLPlatform(ctx context.Context, p *dataConnParam) (*sqlDataPlatform, error) {
	// Validate the connection parameters.
	if err := validateFSConnectionParams(p); err != nil {
		return nil, err
	}

	dbType := p.connectionParams[1]
	projectId := p.connectionParams[0]
	instanceName := p.connectionParams[2]
	dbName := p.connectionParams[3]
	tableName := p.connectionParams[4]

	connectionString := fmt.Sprintf("user@cloudsql(%s:%s)/%s", projectId, instanceName, dbName)

	log.Println(connectionString)

	db, err := sql.Open(dbType, connectionString)

	return &sqlDataPlatform{
		client: db,

		// Join the Cloud SQL path from the parsed parameters.
		query: fmt.Sprintf("select * from %s", tableName),
	}, err

}

// validateConnectionParams is a basic len check of the parameters.
// TODO: Do additional parsing to validate parameters.
func validateSQLConnectionParams(p *dataConnParam) error {

	if len(p.connectionParams) < 5 {
		return errors.New("the url path must be in the form https://host/cloudsql/project/DBType/instance/DBName/table")
	}
	return nil
}
