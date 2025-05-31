package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"

	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type Person struct {
	Name          string
	ID            uint
	Kind          string
	Age           uint8
	Deleted       bool
	Description   string
	Nationality   string
	Cloned        bool
	ClonedFromRef uint
}

func (p *Person) GetEditableFields() []string {
	// Get the editable fields
	return []string{"name", "description", "age", "Nationality"}
}

func (p *Person) PostEditables(editMap map[string]interface{}) {
	// Post the editable fields
	for key, value := range editMap {
		switch key {
		case "name":
			p.Name = value.(string)
		case "description":
			p.Description = value.(string)
		case "age":
			p.Age = value.(uint8)
		case "Nationality":
			p.Nationality = value.(string)
		}
	}
}

func (p *Person) Clone() xBase.Base {
	// Clone the Person
	randomUUID := uuid.New().ID()
	newPerson := &Person{
		Name:          p.Name,
		ID:            uint(randomUUID),
		Kind:          p.Kind,
		Age:           p.Age,
		Deleted:       false,
		Description:   p.Description,
		Nationality:   p.Nationality,
		Cloned:        true,
		ClonedFromRef: p.ID,
	}
	return newPerson
}

func (p *Person) Create(objMap map[string]interface{}) error {
	// Set default values
	p.ID = uint(uuid.New().ID())
	p.Kind = "person"
	p.Deleted = false
	p.Cloned = false
	p.ClonedFromRef = 0

	for key, value := range objMap {
		editableFields := p.GetEditableFields()
		if contains(editableFields, strings.ToLower(key)) {
			field := reflect.ValueOf(p).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

func (p *Person) Update(objMap map[string]interface{}) error {
	// Update the person
	for key, value := range objMap {
		editableFields := p.GetEditableFields()
		if contains(editableFields, strings.ToLower(key)) {
			field := reflect.ValueOf(p).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func (p *Person) Delete() error {
	p.Deleted = true
	err := p.Save()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) Save() error {
	// Save the person
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	// db add/update operations
	return nil
}

func (p *Person) ToString() string {
	// Convert the base to a string
	data := ""
	data += fmt.Sprintf("Name: %s ", p.Name)
	data += fmt.Sprintf("ID: %d ", p.ID)
	data += fmt.Sprintf("Kind: %s ", p.Kind)
	data += fmt.Sprintf("Age: %d ", p.Age)
	data += fmt.Sprintf("Deleted: %t ", p.Deleted)
	data += fmt.Sprintf("Description: %s ", p.Description)
	data += fmt.Sprintf("Nationality: %s", p.Nationality)
	data += fmt.Sprintf("Cloned: %t ", p.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %d ", p.ClonedFromRef)
	return data
}

func (p *Person) ToStatus() map[string]interface{} {
	// Convert the base to a status
	return map[string]interface{}{
		"id":   p.ID,
		"kind": p.Kind,
		"data": map[string]interface{}{
			"name":            p.Name,
			"age":             p.Age,
			"deleted":         p.Deleted,
			"description":     p.Description,
			"Nationality":     p.Nationality,
			"cloned":          p.Cloned,
			"cloned_from_ref": p.ClonedFromRef,
		},
	}
}
