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
	ID            uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `gorm:"column:name;type:varchar(255);not null"`
	Kind          string `gorm:"column:kind;type:varchar(100);not null"`
	Age           uint8  `gorm:"column:age;type:tinyint;not null"`
	Description   string `gorm:"column:description;type:text"`
	Breed         string `gorm:"column:breed;type:varchar(255)"`
	Deleted       bool   `gorm:"column:deleted;type:boolean;default:false"`
	Cloned        bool   `gorm:"column:cloned;type:boolean;default:false"`
	ClonedFromRef uint   `gorm:"column:cloned_from_ref;type:bigint;default:0"`
}

func (a *Animal) GetEditableFields() string {
	// Return the editable fields
	EDITABLE_FIELDS := []string{"name", "description", "age", "breed"}
	editableFields := strings.Join(EDITABLE_FIELDS, "|")
	return editableFields
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
	a.Kind = "person"
	a.Deleted = false
	a.Cloned = false
	a.ClonedFromRef = 0

	for key, value := range objMap {
		if strings.Contains(a.GetEditableFields(), strings.ToLower(key)) {
			field := reflect.ValueOf(a).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return a.Save(nil)
}

func (a *Animal) Update(objMap map[string]interface{}) error {
	// Update the person
	updates := make(map[string]interface{})
	for key, value := range objMap {
		if strings.Contains(a.GetEditableFields(), strings.ToLower(key)) {
			field := reflect.ValueOf(a).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				updates[key] = value
			}
		} else {
			err := fmt.Errorf("field '%s' is not editable", key)
			return xError.NewEditError(err)
		}
	}
	return a.Save(updates)
}

func (a *Animal) Delete() error {
	updates := map[string]interface{}{
		"Deleted": true,
	}
	return a.Save(updates)
}

func (a *Animal) Save(updates map[string]interface{}) error {
	session := xDb.GetSession()
	if session == nil {
		err := errors.New("session not found")
		return xError.NewObjectNotFoundError(err)
	}

	// db add/update operations
	existingRecord := &Person{}
	err := session.ReadRecords(&Person{}, map[string]interface{}{"id": a.ID}, existingRecord)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "record not found") {
			// Record does not exist, create a new one
			err = session.CreateRecords(&Person{}, []interface{}{a})
			if err != nil {
				return xError.NewDBError(err)
			}
		} else {
			// Some other error occurred
			return xError.NewDBError(err)
		}
	} else {
		// Record exists, update it
		err = session.UpdateRecords(&Person{}, map[string]interface{}{"id": a.ID}, updates)
		if err != nil {
			return xError.NewDBError(err)
		}
	}
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
		"id":              a.ID,
		"kind":            a.Kind,
		"name":            a.Name,
		"age":             a.Age,
		"deleted":         a.Deleted,
		"description":     a.Description,
		"Breed":           a.Breed,
		"cloned":          a.Cloned,
		"cloned_from_ref": a.ClonedFromRef,
	}
}

func GetAnimalByID(animalID string) (*Animal, error) {
	session := xDb.GetSession()
	if session == nil {
		err := errors.New("session not found")
		return nil, xError.NewObjectNotFoundError(err)
	}

	animal := &Animal{}
	err := session.ReadRecords(&Person{}, map[string]interface{}{"id": animalID}, animal)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return animal, nil
}
