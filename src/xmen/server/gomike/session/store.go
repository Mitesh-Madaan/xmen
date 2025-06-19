// CRUD operations

package session

import (
	"encoding/json"
	"fmt"

	xError "gomike/error"
)

type Storable interface {
	ToString() string
	ToStatus() map[string]interface{}
	Init()
}

func CreateRecord[T Storable](objDetails []byte) (obj *T, err error) {
	// Create a new record in the database
	// This function will parse the objDetails and create a new record
	// using the appropriate model based on the details provided.

	obj = new(T)
	(*obj).Init() // Initialize the object
	fmt.Printf("Creating record with details: %s\n", string(objDetails))
	err = json.Unmarshal(objDetails, obj)
	if err != nil {
		return nil, xError.NewParseError(err)
	}

	err = dbSession.CreateRecords(&obj, []interface{}{obj})
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return obj, nil
}

func ReadRecord[T Storable](objID string) (obj *T, err error) {
	// Read a record from the database by ID
	// This function will retrieve the record based on the ID provided.

	obj = new(T)
	err = dbSession.ReadRecords(&obj, map[string]interface{}{"id": objID}, obj)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return obj, nil
}

// UpdateRecord()

// DeleteRecord()
