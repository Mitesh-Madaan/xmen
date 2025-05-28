package models

import (
	"encoding/json"
	"errors"
	"fmt"
	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type Person struct {
	base *xBase.Base
}

func NewPerson(b *xBase.Base) IBase {
	// Create a new person
	return &Person{base: b}
}

func (p *Person) Delete() error {
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	p.base.MarkDeleted()
	err := p.Save()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) Clone() IBase {
	// Clone the Person
	b := p.base.Clone()
	clone := NewPerson(b)
	return clone
}

func (p *Person) Edit(config map[string]interface{}) error {
	// Edit the person
	// editableFields := p.base.GetEditableFields()
	editMap := make(map[string]interface{})

	for key, value := range config {
		// if !xBase.Contains(editableFields, key) {
		// 	return xError.NewValidationError("invalid field")
		// }
		editMap[key] = value
	}

	p.base.PostEditables(editMap)
	err := p.Save()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) PreviewValidate() bool {
	// Validate the
	if p.base == nil {
		return false
	}
	return p.base.PreviewValidate()
}

func (p *Person) Save() error {
	// Save the person
	check := p.PreviewValidate()
	if !check {
		err := errors.New(p.base.GetMessages())
		return xError.NewValidationError(err)
	}

	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	dbIdentity := p.base.GetDBIdentifier()
	_, err := directory.Get(dbIdentity["id"], dbIdentity["kind"])
	details := p.ToJson()
	if details == nil {
		err := errors.New("failed to convert person to JSON")
		return xError.NewValidationError(err)
	}
	// fmt.Printf("details of person: %v\n", details)
	if err != nil && err.Error() == "record not found" {
		directory.Add(dbIdentity["id"], dbIdentity["kind"], details)
	} else {
		directory.Update(dbIdentity["id"], dbIdentity["kind"], details)
	}
	return nil
}

func (p *Person) Fill(details []byte) error {
	// Fill the person
	return p.FromJson(details)
}

func (p *Person) ToString() string {
	// Convert the person to a string
	return fmt.Sprintf("Person: %s\n", p.base.ToString())
}

func (p *Person) ToStatus() map[string]interface{} {
	// Convert the person to a status
	return p.base.ToStatus()
}

func (p *Person) FromJson(details []byte) error {
	// Fill the person from JSON
	err := json.Unmarshal(details, &p.base)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON: %v\n", err)
		return err
	}
	return nil
}

func (p *Person) ToJson() []byte {
	// Convert the person to JSON
	details, err := p.base.ToJson()
	if err != nil {
		fmt.Printf("Error converting to JSON: %v\n", err)
		return nil
	}
	return details
}

func (p *Person) GetBase() *xBase.Base {
	// Get the base
	return p.base
}
