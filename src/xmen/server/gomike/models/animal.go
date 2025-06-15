package models

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type Animal struct {
	ID            string `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;type:varchar(255);not null"`
	Kind          string `gorm:"column:kind;type:varchar(100);not null"`
	Age           int    `gorm:"column:age;type:int;not null"`
	Description   string `gorm:"column:description;type:text"`
	Breed         string `gorm:"column:breed;type:varchar(255)"`
	Cloned        bool   `gorm:"column:cloned;type:boolean;default:false"`
	ClonedFromRef string `gorm:"column:cloned_from_ref;type:varchar(100);default:''"`
}

func (a *Animal) Clone() xBase.Base {
	// Clone the Animal
	randomUUID := uuid.New().String()
	newAnimal := &Animal{
		Name:          a.Name,
		ID:            randomUUID,
		Kind:          a.Kind,
		Age:           a.Age,
		Description:   a.Description,
		Breed:         a.Breed,
		Cloned:        true,
		ClonedFromRef: a.ID,
	}
	return newAnimal
}

func (a *Animal) Create(dbSession *xDb.DBSession, objDetails []byte) error {
	// Parse the JSON object details
	err := json.Unmarshal(objDetails, a)
	if err != nil {
		return xError.NewParseError(err)
	}

	// Setting default values
	a.Kind = "animal"
	a.Cloned = false
	a.ClonedFromRef = ""

	fmt.Println("Creating animal with details:", a.ToString())
	// Create the person
	err = dbSession.CreateRecords(&Animal{}, []interface{}{a})
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (a *Animal) Update(dbSession *xDb.DBSession, objDetails []byte) error {
	// Update the editable fields
	if objDetails != nil {
		err := json.Unmarshal(objDetails, &a)
		if err != nil {
			return xError.NewParseError(err)
		}
	}
	fmt.Println("Updating animal with details:", a.ToString())
	// Update the person
	err := dbSession.UpdateRecords(&Animal{}, map[string]interface{}{"id": a.ID}, a.ToStatus())
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
	data += fmt.Sprintf("ID: %s ", a.ID)
	data += fmt.Sprintf("Kind: %s ", a.Kind)
	data += fmt.Sprintf("Age: %d ", a.Age)
	data += fmt.Sprintf("Description: %s ", a.Description)
	data += fmt.Sprintf("Breed: %s ", a.Breed)
	data += fmt.Sprintf("Cloned: %t ", a.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %s ", a.ClonedFromRef)
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

func GetAnimalByID(dbSession *xDb.DBSession, animalID string) (*Animal, error) {
	animal := &Animal{}
	err := dbSession.ReadRecords(&Animal{}, map[string]interface{}{"id": animalID}, animal)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return animal, nil
}
