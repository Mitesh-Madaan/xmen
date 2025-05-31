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

type Animal struct {
	Name          string
	ID            uint
	Kind          string
	Age           uint8
	Deleted       bool
	Description   string
	Breed         string
	Cloned        bool
	ClonedFromRef uint
}

func (a *Animal) GetEditableFields() []string {
	// Get the editable fields
	return []string{"name", "description", "age", "Breed"}
}

func (a *Animal) PostEditables(editMap map[string]interface{}) {
	// Post the editable fields
	for key, value := range editMap {
		switch key {
		case "name":
			a.Name = value.(string)
		case "description":
			a.Description = value.(string)
		case "age":
			a.Age = value.(uint8)
		case "Breed":
			a.Breed = value.(string)
		}
	}
}

func (a *Animal) Clone() xBase.Base {
	// Clone the Animal
	randomUUID := uuid.New().ID()
	newAnimal := &Animal{
		Name:          a.Name,
		ID:            uint(randomUUID),
		Kind:          a.Kind,
		Age:           a.Age,
		Deleted:       false,
		Description:   a.Description,
		Breed:         a.Breed,
		Cloned:        true,
		ClonedFromRef: a.ID,
	}
	return newAnimal
}

func (a *Animal) Create(objMap map[string]interface{}) error {
	// Set default values
	a.ID = uint(uuid.New().ID())
	a.Kind = "animal"
	a.Deleted = false
	a.Cloned = false
	a.ClonedFromRef = 0

	for key, value := range objMap {
		editableFields := a.GetEditableFields()
		if contains(editableFields, strings.ToLower(key)) {
			field := reflect.ValueOf(a).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}

	return nil
}

func (a *Animal) Update(objMap map[string]interface{}) error {
	// Update the animal
	for key, value := range objMap {
		editableFields := a.GetEditableFields()
		if contains(editableFields, strings.ToLower(key)) {
			field := reflect.ValueOf(a).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}

func (a *Animal) Delete() error {
	a.Deleted = true
	err := a.Save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Animal) Save() error {
	// Save the animal
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	// db add/update operations
	return nil
}

func (a *Animal) ToString() string {
	// Convert the base to a string
	data := ""
	data += fmt.Sprintf("Name: %s ", a.Name)
	data += fmt.Sprintf("ID: %d ", a.ID)
	data += fmt.Sprintf("Kind: %s ", a.Kind)
	data += fmt.Sprintf("Age: %d ", a.Age)
	data += fmt.Sprintf("Deleted: %t ", a.Deleted)
	data += fmt.Sprintf("Description: %s ", a.Description)
	data += fmt.Sprintf("Breed: %s", a.Breed)
	data += fmt.Sprintf("Cloned: %t ", a.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %d ", a.ClonedFromRef)
	return data
}

func (a *Animal) ToStatus() map[string]interface{} {
	// Convert the base to a status
	return map[string]interface{}{
		"id":   a.ID,
		"kind": a.Kind,
		"data": map[string]interface{}{
			"name":            a.Name,
			"age":             a.Age,
			"deleted":         a.Deleted,
			"description":     a.Description,
			"Breed":           a.Breed,
			"cloned":          a.Cloned,
			"cloned_from_ref": a.ClonedFromRef,
		},
	}
}
