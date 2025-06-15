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
	err = dbSession.CreateRecords(&Person{}, []interface{}{p})
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
	err := dbSession.UpdateRecords(&Person{}, map[string]interface{}{"id": p.ID}, p.ToStatus())
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (p *Person) Delete(dbSession *xDb.DBSession) error {
	err := dbSession.DeleteRecords(&Person{}, map[string]interface{}{"id": p.ID})
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
	data := ""
	data += fmt.Sprintf("Name: %s ", p.Name)
	data += fmt.Sprintf("ID: %s ", p.ID)
	data += fmt.Sprintf("Kind: %s ", p.Kind)
	data += fmt.Sprintf("Age: %d ", p.Age)
	data += fmt.Sprintf("Description: %s ", p.Description)
	data += fmt.Sprintf("Nationality: %s ", p.Nationality)
	data += fmt.Sprintf("Cloned: %t ", p.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %s ", p.ClonedFromRef)
	return data
}

func (p *Person) ToStatus() map[string]interface{} {
	// Convert the base to a status
	return map[string]interface{}{
		"id":              p.ID,
		"kind":            p.Kind,
		"name":            p.Name,
		"age":             p.Age,
		"description":     p.Description,
		"Nationality":     p.Nationality,
		"cloned":          p.Cloned,
		"cloned_from_ref": p.ClonedFromRef,
	}
}

func GetPersonByID(dbSession *xDb.DBSession, personID string) (*Person, error) {
	person := &Person{}
	err := dbSession.ReadRecords(&Person{}, map[string]interface{}{"id": personID}, person)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return person, nil
}
