package models

import (
	"encoding/json"
	"errors"
	"fmt"
	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type Animal struct {
	base *xBase.Base
}

func NewAnimal(b *xBase.Base) IBase {
	// Create a new animal
	return &Animal{base: b}
}

func (a *Animal) GetBase() *xBase.Base {
	// Get the base
	return a.base
}

func (a *Animal) Delete() error {
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	a.base.MarkDeleted()
	err := a.Save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Animal) Clone() IBase {
	// Clone the Animal
	b := a.base.Clone()
	clone := NewAnimal(b)
	return clone
}

func (a *Animal) Edit(config map[string]interface{}) error {
	// Edit the animal
	// editableFields := a.base.GetEditableFields()
	editMap := make(map[string]interface{})

	for key, value := range config {
		// if !xBase.Contains(editableFields, key) {
		// 	return xError.NewValidationError("invalid field")
		// }
		editMap[key] = value
	}

	a.base.PostEditables(editMap)
	err := a.Save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Animal) PreviewValidate() bool {
	// Validate the
	if a.base == nil {
		return false
	}
	return a.base.PreviewValidate()
}

func (a *Animal) Save() error {
	// Save the animal
	check := a.PreviewValidate()
	if !check {
		err := errors.New(a.base.GetMessages())
		return xError.NewValidationError(err)
	}

	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	dbIdentity := a.base.GetDBIdentifier()
	_, err := directory.Get(dbIdentity["id"], dbIdentity["kind"])
	if err != nil && err.Error() == "record not found" {
		directory.Add(dbIdentity["id"], dbIdentity["kind"], a.ToJson())
	} else {
		directory.Update(dbIdentity["id"], dbIdentity["kind"], a.ToJson())
	}
	return nil
}

func (a *Animal) Fill(details []byte) error {
	// Fill the animal
	return a.FromJson(details)
}

func (a *Animal) ToString() string {
	// Convert the animal to a string
	return fmt.Sprintf("Animal: %s\n", a.base.ToString())
}

func (a *Animal) ToStatus() map[string]interface{} {
	// Convert the animal to a status
	return a.base.ToStatus()
}

func (a *Animal) FromJson(details []byte) error {
	// Convert the details to an Animal
	err := json.Unmarshal(details, &a.base)
	if err != nil {
		fmt.Printf("failed to unmarshal details: %v", err)
		return err
	}
	return nil
}

func (a *Animal) ToJson() []byte {
	// Convert the animal to JSON
	details, err := a.base.ToJson()
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return nil
	}
	return details
}
