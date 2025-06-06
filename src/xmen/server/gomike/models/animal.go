package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type Animal struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `gorm:"column:name;type:varchar(255);not null"`
	Kind          string `gorm:"column:kind;type:varchar(100);not null"`
	Age           int    `gorm:"column:age;type:int;not null"`
	Description   string `gorm:"column:description;type:text"`
	Breed         string `gorm:"column:breed;type:varchar(255)"`
	Cloned        bool   `gorm:"column:cloned;type:boolean;default:false"`
	ClonedFromRef uint64 `gorm:"column:cloned_from_ref;type:bigint;default:0"`
}

func (a *Animal) PostEditableFields(objMap map[string]interface{}) error {
	// Update the editable fields
	for key, value := range objMap {
		// Print the key and value for debugging
		fmt.Printf("Key: %s, Value: %v\n", key, value)
		switch strings.ToLower(key) {
		case "name":
			a.Name = fmt.Sprintf("%v", value)
		case "age":
			a.Age = int(value.(float64)) // Assuming value is a float64, adjust as necessary
		case "description":
			a.Description = fmt.Sprintf("%v", value)
		case "breed":
			a.Breed = fmt.Sprintf("%v", value)
		default:
			err := fmt.Errorf("field '%s' is not editable", key)
			return xError.NewEditError(err)
		}
	}
	return nil
}

func (a *Animal) Clone() xBase.Base {
	// Clone the Animal
	randomUUID := uuid.New().ID()
	newAnimal := &Animal{
		Name:          a.Name,
		ID:            uint64(randomUUID),
		Kind:          a.Kind,
		Age:           a.Age,
		Description:   a.Description,
		Breed:         a.Breed,
		Cloned:        true,
		ClonedFromRef: a.ID,
	}
	return newAnimal
}

func (a *Animal) Create(dbSession *xDb.DBSession, objMap map[string]interface{}) error {
	// Set default values
	a.Kind = "animal"
	a.Cloned = false
	a.ClonedFromRef = 0

	// Update the editable fields
	if objMap != nil {
		if objMap["id"] != nil {
			// If an ID is provided, set it
			a.ID = objMap["id"].(uint64)
			delete(objMap, "id") // Remove ID from objMap to avoid conflicts
		} else {
			// If no ID is provided, generate a new one
			a.ID = uint64(uuid.New().ID())
		}
		err := a.PostEditableFields(objMap)
		if err != nil {
			return err
		}
	}
	fmt.Println("Creating animal with details:", a.ToString())
	// Create the animal
	err := dbSession.CreateRecords(&Animal{}, []interface{}{a})
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (a *Animal) Update(dbSession *xDb.DBSession, editMap map[string]interface{}) error {
	// Update the editable fields
	if editMap != nil {
		err := a.PostEditableFields(editMap)
		if err != nil {
			return err
		}
	}
	fmt.Println("Updating animal with details:", a.ToString())
	// Update the animal
	err := dbSession.UpdateRecords(&Animal{}, map[string]interface{}{"id": a.ID}, editMap)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (a *Animal) Delete(dbSession *xDb.DBSession) error {
	err := dbSession.DeleteRecords(&Animal{}, map[string]interface{}{"id": a.ID})
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (a *Animal) Save(dbSession *xDb.DBSession, updates map[string]interface{}) error {
	// No need to implement this method for Animal as it is not used in the current context
	return nil
}

func (a *Animal) ToString() string {
	// Convert the base to a string
	data := ""
	data += fmt.Sprintf("Name: %s ", a.Name)
	data += fmt.Sprintf("ID: %d ", a.ID)
	data += fmt.Sprintf("Kind: %s ", a.Kind)
	data += fmt.Sprintf("Age: %d ", a.Age)
	data += fmt.Sprintf("Description: %s ", a.Description)
	data += fmt.Sprintf("Breed: %s ", a.Breed)
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
		"description":     a.Description,
		"Breed":           a.Breed,
		"cloned":          a.Cloned,
		"cloned_from_ref": a.ClonedFromRef,
	}
}

func GetAnimalByID(dbSession *xDb.DBSession, animalID uint64) (*Animal, error) {
	animal := &Animal{}
	err := dbSession.ReadRecords(&Animal{}, map[string]interface{}{"id": animalID}, animal)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return animal, nil
}
