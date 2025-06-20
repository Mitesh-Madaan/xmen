package models

import (
	"encoding/json"
	"fmt"

	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"

	"github.com/google/uuid"
)

type Person struct {
	ID            string `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name;type:varchar(100);not null"`
	Kind          string `gorm:"column:kind;type:varchar(50);not null"`
	Age           int    `gorm:"column:age;type:int"`
	Description   string `gorm:"column:description;type:text"`
	Nationality   string `gorm:"column:nationality;type:varchar(100)"`
	Cloned        bool   `gorm:"column:cloned;default:false"`
	ClonedFromRef string `gorm:"column:cloned_from_ref;type:varchar(100);default:''"`
}

func (p *Person) Clone() xBase.Base {
	// Clone the Person
	randomUUID := uuid.New().String()
	newPerson := &Person{
		Name:          p.Name,
		ID:            randomUUID,
		Kind:          p.Kind,
		Age:           p.Age,
		Description:   p.Description,
		Nationality:   p.Nationality,
		Cloned:        true,
		ClonedFromRef: p.ID,
	}
	return newPerson
}

func (p *Person) Create(dbSession *xDb.DBSession, objDetails []byte) error {
	// Parse the JSON object details
	err := json.Unmarshal(objDetails, p)
	if err != nil {
		return xError.NewParseError(err)
	}

	// Setting default values
	p.Kind = "person"
	p.Cloned = false
	p.ClonedFromRef = ""

	fmt.Println("Creating person with details:", p.ToString())
	// Create the person
	err = dbSession.CreateRecord(p)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (p *Person) Update(dbSession *xDb.DBSession, objDetails []byte) error {
	// Update the editable fields
	if objDetails != nil {
		err := json.Unmarshal(objDetails, &p)
		if err != nil {
			return xError.NewParseError(err)
		}
	}
	fmt.Println("Updating person with details:", p.ToString())
	// Update the person
	err := dbSession.UpdateRecord(p)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (p *Person) Delete(dbSession *xDb.DBSession) error {
	err := dbSession.DeleteRecord(p)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (p *Person) Save(dbSession *xDb.DBSession, updates map[string]interface{}) error {
	// db add/update operations
	return nil // No need to implement this method for Person as Create and Update handle it
}

func (p *Person) ToString() string {
	// Convert the base to a string
	data := fmt.Sprintf(`
		Name: %s
		ID: %s
		Kind: %s
		Age: %d 
		Description: %s
		Nationality: %s
		Cloned: %t 
		Cloned From Ref: %s
		`, p.Name, p.ID, p.Kind, p.Age, p.Description, p.Nationality, p.Cloned, p.ClonedFromRef)
	return data
}

func GetPersonByID(dbSession *xDb.DBSession, personID string) (*Person, error) {
	person := &Person{}
	err := dbSession.ReadRecord(map[string]interface{}{"id": personID}, person)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return person, nil
}
