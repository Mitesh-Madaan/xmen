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
	ID            uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `gorm:"column:name;type:varchar(100);not null"`
	Kind          string `gorm:"column:kind;type:varchar(50);not null"`
	Age           int    `gorm:"column:age;type:int;not null"`
	Description   string `gorm:"column:description;type:text"`
	Nationality   string `gorm:"column:nationality;type:varchar(100)"`
	Deleted       bool   `gorm:"column:deleted;default:false"`
	Cloned        bool   `gorm:"column:cloned;default:false"`
	ClonedFromRef uint   `gorm:"column:cloned_from_ref;default:0"`
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

func (p *Person) Create(dbSession *xDb.DBSession, objMap map[string]interface{}) error {
	// Set default values
	p.ID = uint(uuid.New().ID())
	p.Kind = "person"
	p.Deleted = false
	p.Cloned = false
	p.ClonedFromRef = 0

	// Update the editable fields
	if objMap != nil {
		err := p.PostEditableFields(objMap)
		if err != nil {
			return err
		}
	}
	fmt.Println("Creating person with details:", p.ToString())
	// Save the person
	return p.Save(dbSession, nil)
}

func (p *Person) Update(dbSession *xDb.DBSession, editMap map[string]interface{}) error {
	// Update the editable fields
	if editMap != nil {
		err := p.PostEditableFields(editMap)
		if err != nil {
			return err
		}
	}
	// Save the person with updates
	return p.Save(dbSession, editMap)
}

func (p *Person) Delete(dbSession *xDb.DBSession) error {
	updates := map[string]interface{}{
		"Deleted": true,
	}
	return p.Save(dbSession, updates)
}

func (p *Person) Save(dbSession *xDb.DBSession, updates map[string]interface{}) error {
	// db add/update operations
	existingRecord := &Person{}
	err := dbSession.ReadRecords(&Person{}, map[string]interface{}{"id": p.ID}, existingRecord)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "record not found") {
			// Record does not exist, create a new one
			err = dbSession.CreateRecords(&Person{}, []interface{}{p})
			if err != nil {
				return xError.NewDBError(err)
			}
		} else {
			// Some other error occurred
			return xError.NewDBError(err)
		}
	} else {
		// Record exists, update it
		err = dbSession.UpdateRecords(&Person{}, map[string]interface{}{"id": p.ID}, updates)
		if err != nil {
			return xError.NewDBError(err)
		}
	}
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
	data += fmt.Sprintf("Nationality: %s ", p.Nationality)
	data += fmt.Sprintf("Cloned: %t ", p.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %d ", p.ClonedFromRef)
	return data
}

func (p *Person) ToStatus() map[string]interface{} {
	// Convert the base to a status
	return map[string]interface{}{
		"id":              p.ID,
		"kind":            p.Kind,
		"name":            p.Name,
		"age":             p.Age,
		"deleted":         p.Deleted,
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
