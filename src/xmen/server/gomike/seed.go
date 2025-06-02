package main

import (
	"fmt"

	"lib/dbchef"

	xError "gomike/error"
	xModels "gomike/models"
)

var session *dbchef.DBSession

// SeedTables creates and seeds Person and Animal tables
func SeedTables() error {
	models := []interface{}{
		&xModels.Person{},
		&xModels.Animal{},
	}

	err := session.SeedTables(models)
	if err != nil {
		err := fmt.Errorf("failed to seed tables: %w", err)
		return xError.NewDBError(err)
	}
	return nil
}

func SeedRecords() error {
	// Seed the directory
	var err error
	session = dbchef.GetSession()
	if session == nil {
		err = fmt.Errorf("session not found")
		return xError.NewObjectNotFoundError(err)
	}

	// data, err := ioutil.ReadFile("/Users/mitesh.madaan/xmen/records.json")
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return err
	// }

	return nil
}

// func StoreRecords() error {
// 	// Store the directory
// 	// var err error
// 	directory = dbchef.GetDirectory()

// 	var storeRecords []*base.Base
// 	for _, record := range directory.Records {
// 		obj, err := xModels.GetObjectFromDB(record.Id, record.Kind)
// 		if err != nil {
// 			fmt.Printf("Error: %v\n", err)
// 		}
// 		fmt.Printf("Store Record: %s\n", obj.ToString())
// 		storeRecords = append(storeRecords, obj.GetBase())
// 	}

// 	data, err := json.MarshalIndent(storeRecords, "", "  ")
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return err
// 	}

// 	if err := ioutil.WriteFile("/Users/mitesh.madaan/xmen/output.json", data, os.ModePerm); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return err
// 	}
// 	return nil
// }

// func PrintRecords() error {
// 	// Print the directory
// 	directory = dbchef.GetDirectory()
// 	for _, record := range directory.Records {
// 		fmt.Printf("Record ID: %s, Kind: %s, details: %v\n", record.Id, record.Kind, record.ObjDetails)
// 	}
// 	return nil
// }
