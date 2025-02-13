package models

import (
	"errors"
	"fmt"
	xError "gomike/error"

	"github.com/google/uuid"

	xBase "lib/base"
	xDb "lib/dbchef"
	xPb "lib/pb"

	"google.golang.org/protobuf/proto"
)

type Person struct {
	base *xBase.Base
}

func NewPerson(b *xBase.Base) IBase {
	// Create a new person
	return &Person{base: b}
}

func (p *Person) Delete() error {
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	p.base.MarkDeleted()
	err := p.Save()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) Clone() IBase {
	// Clone the Person
	b := p.base.Clone()
	clone := NewPerson(b)
	return clone
}

func (p *Person) Edit(config map[string]interface{}) error {
	// Edit the person
	// editableFields := p.base.GetEditableFields()
	editMap := make(map[string]interface{})

	for key, value := range config {
		// if !xBase.Contains(editableFields, key) {
		// 	return xError.NewValidationError("invalid field")
		// }
		editMap[key] = value
	}

	p.base.PostEditables(editMap)
	err := p.Save()
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) PreviewValidate() bool {
	// Validate the
	if p.base == nil {
		return false
	}
	return p.base.PreviewValidate()
}

func (p *Person) Save() error {
	// Save the person
	check := p.PreviewValidate()
	if !check {
		err := errors.New(p.base.GetMessages())
		return xError.NewValidationError(err)
	}

	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	dbIdentity := p.base.GetDBIdentifier()
	_, err := directory.Get(dbIdentity["id"], dbIdentity["kind"])
	if err != nil && err.Error() == "record not found" {
		directory.Add(dbIdentity["id"], dbIdentity["kind"], p.ToProto())
	} else {
		directory.Update(dbIdentity["id"], dbIdentity["kind"], p.ToProto())
	}
	return nil
}

func (p *Person) Fill(msg proto.Message) error {
	// Fill the person
	return p.FromProto(msg.(*xPb.Person))
}

func (p *Person) ToString() string {
	// Convert the person to a string
	return fmt.Sprintf("Person: %s\n", p.base.ToString())
}

func (p *Person) ToStatus() map[string]interface{} {
	// Convert the person to a status
	return p.base.ToStatus()
}

func (p *Person) ToProto() *xPb.Person {
	// Convert the person to a proto
	msg := &xPb.Person{}
	dbIdentity := p.base.GetDBIdentifier()
	id := dbIdentity["id"]
	msg.Id = &id
	kind := dbIdentity["kind"]
	msg.Kind = &kind
	name := p.base.Name
	msg.Name = &name
	age := int64(p.base.Age)
	msg.Age = &age
	deleted := p.base.Deleted
	msg.Deleted = &deleted
	description := p.base.Description
	msg.Description = &description
	cloned := p.base.Cloned
	msg.Cloned = &cloned
	clonedFromRef := p.base.ClonedFromRef.String()
	msg.ClonedFromRef = &clonedFromRef
	return msg
}

func (p *Person) FromProto(msg *xPb.Person) error {
	// Convert the person from a proto
	id, err := uuid.Parse(msg.GetId())
	if err != nil {
		return err
	}
	p.base.Name = msg.GetName()
	p.base.ID = id
	p.base.Kind = msg.GetKind()
	p.base.Age = int(msg.GetAge())
	p.base.Deleted = msg.GetDeleted()
	p.base.Description = msg.GetDescription()
	p.base.Cloned = msg.GetCloned()
	clonedFromRef, err := uuid.Parse(msg.GetClonedFromRef())
	if err != nil {
		return err
	}
	p.base.ClonedFromRef = clonedFromRef
	return nil
}

func (p *Person) GetBase() *xBase.Base {
	// Get the base
	return p.base
}
