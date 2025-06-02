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
	ID            uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Name          string `gorm:"column:name;type:varchar(100);not null"`
	Kind          string `gorm:"column:kind;type:varchar(50);not null"`
	Age           uint8  `gorm:"column:age;not null"`
	Description   string `gorm:"column:description;type:text"`
	Nationality   string `gorm:"column:nationality;type:varchar(100)"`
	Deleted       bool   `gorm:"column:deleted;default:false"`
	Cloned        bool   `gorm:"column:cloned;default:false"`
	ClonedFromRef uint   `gorm:"column:cloned_from_ref;default:0"`
}

func (p *Person) GetEditableFields() string {
	// Return the editable fields
	EDITABLE_FIELDS := []string{"name", "age", "description", "nationality"}
	editableFields := strings.Join(EDITABLE_FIELDS, "|")
	return editableFields
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
		if strings.Contains(p.GetEditableFields(), strings.ToLower(key)) {
			field := reflect.ValueOf(p).Elem().FieldByNameFunc(func(fieldName string) bool {
				return strings.EqualFold(fieldName, key)
			})
			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(value))
			}
		}
	}
	return p.Save(nil)
}

func (p *Person) Update(objMap map[string]interface{}) error {
	// Update the person
	updates := make(map[string]interface{})
	for key, value := range objMap {
		if strings.Contains(p.GetEditableFields(), strings.ToLower(key)) {
			field := reflect.ValueOf(p).Elem().FieldByNameFunc(func(fieldName string) bool {
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
	return p.Save(updates)
}

func (p *Person) Delete() error {
	updates := map[string]interface{}{
		"Deleted": true,
	}
	return p.Save(updates)
}

func (p *Person) Save(updates map[string]interface{}) error {
	session := xDb.GetSession()
	if session == nil {
		err := errors.New("session not found")
		return xError.NewObjectNotFoundError(err)
	}

	// db add/update operations
	existingRecord := &Person{}
	err := session.ReadRecords(&Person{}, map[string]interface{}{"id": p.ID}, existingRecord)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "record not found") {
			// Record does not exist, create a new one
			err = session.CreateRecords(&Person{}, []interface{}{p})
			if err != nil {
				return xError.NewDBError(err)
			}
		} else {
			// Some other error occurred
			return xError.NewDBError(err)
		}
	} else {
		// Record exists, update it
		err = session.UpdateRecords(&Person{}, map[string]interface{}{"id": p.ID}, updates)
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
	data += fmt.Sprintf("Nationality: %s", p.Nationality)
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

func GetPersonByID(personID string) (*Person, error) {
	session := xDb.GetSession()
	if session == nil {
		err := errors.New("session not found")
		return nil, xError.NewObjectNotFoundError(err)
	}

	person := &Person{}
	err := session.ReadRecords(&Person{}, map[string]interface{}{"id": personID}, person)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return person, nil
}
