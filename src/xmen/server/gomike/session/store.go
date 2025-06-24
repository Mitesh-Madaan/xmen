// CRUD operations

package session

import (
	"encoding/json"
	"fmt"

	xError "gomike/error"
)

type Storable any

func CreateRecord[T Storable](obj T, objDetails []byte) error {
	// Create a new record in the database
	// This function will parse the objDetails and create a new record
	// using the appropriate model based on the details provided.

	fmt.Printf("Creating record with details: %s\n", string(objDetails))
	err := json.Unmarshal(objDetails, obj)
	if err != nil {
		return xError.NewParseError(err)
	}

	err = dbSession.CreateRecord(&obj)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func ReadRecord[T Storable](objID string) (obj *T, err error) {
	// Read a record from the database by ID
	// This function will retrieve the record based on the ID provided.

	obj = new(T)
	err = dbSession.ReadRecord(map[string]interface{}{"id": objID}, &obj)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return obj, nil
}

// UpdateRecord()

// DeleteRecord()
