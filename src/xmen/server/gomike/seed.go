package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"lib/base"
	"lib/dbchef"

	xModels "gomike/models"
)

var directory *dbchef.Directory

func InitDirectory() error {
	directory = dbchef.GetDirectory()
	if directory == nil {
		return fmt.Errorf("directory not found")
	}
	return nil
}

func SeedDirectory() error {
	// Seed the directory
	var err error
	directory = dbchef.GetDirectory()
	data, err := ioutil.ReadFile("/Users/mitesh.madaan/xmen/records.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	var dbRecords []base.Base
	if err := json.Unmarshal(data, &dbRecords); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	for _, record := range dbRecords {
		// fmt.Printf("Base Record seeding: %v\n", record)
		obj := xModels.GetBase(&record)
		// fmt.Printf("Base obj: %s\n", obj.ToString())
		err = obj.Save()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
	}
	return nil
}

func StoreDirectory() error {
	// Store the directory
	// var err error
	directory = dbchef.GetDirectory()

	var storeRecords []*base.Base
	for _, record := range directory.Records {
		obj, err := xModels.GetObjectFromDB(record.Id, record.Kind)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("Store Record: %s\n", obj.ToString())
		storeRecords = append(storeRecords, obj.GetBase())
	}

	data, err := json.MarshalIndent(storeRecords, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	if err := ioutil.WriteFile("/Users/mitesh.madaan/xmen/output.json", data, os.ModePerm); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	return nil
}

func PrintDirectory() error {
	// Print the directory
	directory = dbchef.GetDirectory()
	for _, record := range directory.Records {
		fmt.Printf("Record ID: %s, Kind: %s, details: %v\n", record.Id, record.Kind, record.ObjDetails)
	}
	return nil
}
