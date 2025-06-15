package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	xError "gomike/error"
	xBase "lib/base"
	xDb "lib/dbchef"
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

func (p *Person) PostEditableFields(objMap map[string]interface{}) error {
	// Update the editable fields
	for key, value := range objMap {
		// Print the key and value for debugging
		fmt.Printf("Key: %s, Value: %v\n", key, value)
		switch strings.ToLower(key) {
		case "name":
			p.Name = fmt.Sprintf("%v", value)
		case "age":
			p.Age = int(value.(float64)) // Assuming value is a float64, adjust as necessary
		case "description":
			p.Description = fmt.Sprintf("%v", value)
		case "nationality":
			p.Nationality = fmt.Sprintf("%v", value)
		default:
			err := fmt.Errorf("field '%s' is not editable", key)
			return xError.NewEditError(err)
		}
	}
	return nil
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

func (p *Person) Create(dbSession *xDb.DBSession, objMap map[string]interface{}) error {
	// Set default values
	p.Kind = "person"
	p.Cloned = false
	p.ClonedFromRef = ""

	// Update the editable fields
	if objMap != nil {
		if objMap["id"] != nil {
			// If ID is provided, set it
			p.ID = objMap["id"].(string)
			delete(objMap, "id") // Remove ID from objMap to avoid conflicts
		} else {
			// Generate a new ID if not provided
			p.ID = uuid.New().String()
		}
		err := p.PostEditableFields(objMap)
		if err != nil {
			return err
		}
	}
	fmt.Println("Creating person with details:", p.ToString())
	// Create the person
	err := dbSession.CreateRecords(&Person{}, []interface{}{p})
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func (p *Person) Update(dbSession *xDb.DBSession, editMap map[string]interface{}) error {
	// Update the editable fields
	if editMap != nil {
		err := p.PostEditableFields(editMap)
		if err != nil {
			return err
		}
	}
	fmt.Println("Updating person with details:", p.ToString())
	// Update the person
	err := dbSession.UpdateRecords(&Person{}, map[string]interface{}{"id": p.ID}, editMap)
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
