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
	err = dbSession.CreateRecord(a)
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
	err := dbSession.UpdateRecord(a)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (a *Animal) Delete(dbSession *xDb.DBSession) error {
	err := dbSession.DeleteRecord(a)
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
	data := fmt.Sprintf(`
		Name: %s
		ID: %s
		Kind: %s
		Age: %d 
		Description: %s
		Breed: %s
		Cloned: %t 
		Cloned From Ref: %s
		`, a.Name, a.ID, a.Kind, a.Age, a.Description, a.Breed, a.Cloned, a.ClonedFromRef)
	return data
}

func GetAnimalByID(dbSession *xDb.DBSession, animalID string) (*Animal, error) {
	animal := &Animal{}
	err := dbSession.ReadRecord(map[string]interface{}{"id": animalID}, animal)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return animal, nil
}
