package models

import (
	"errors"
	"fmt"
	xError "gomike/error"

	xBase "lib/base"
	xDb "lib/dbchef"
	xPb "lib/pb"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type Animal struct {
	base *xBase.Base
}

func NewAnimal(b *xBase.Base) IBase {
	// Create a new animal
	return &Animal{base: b}
}

func (a *Animal) GetBase() *xBase.Base {
	// Get the base
	return a.base
}

func (a *Animal) Delete() error {
	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	a.base.MarkDeleted()
	err := a.Save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Animal) Clone() IBase {
	// Clone the Animal
	b := a.base.Clone()
	clone := NewAnimal(b)
	return clone
}

func (a *Animal) Edit(config map[string]interface{}) error {
	// Edit the animal
	// editableFields := a.base.GetEditableFields()
	editMap := make(map[string]interface{})

	for key, value := range config {
		// if !xBase.Contains(editableFields, key) {
		// 	return xError.NewValidationError("invalid field")
		// }
		editMap[key] = value
	}

	a.base.PostEditables(editMap)
	err := a.Save()
	if err != nil {
		return err
	}
	return nil
}

func (a *Animal) PreviewValidate() bool {
	// Validate the
	if a.base == nil {
		return false
	}
	return a.base.PreviewValidate()
}

func (a *Animal) Save() error {
	// Save the animal
	check := a.PreviewValidate()
	if !check {
		err := errors.New(a.base.GetMessages())
		return xError.NewValidationError(err)
	}

	directory := xDb.GetDirectory()
	if directory == nil {
		err := errors.New("directory not found")
		return xError.NewObjectNotFoundError(err)
	}

	dbIdentity := a.base.GetDBIdentifier()
	_, err := directory.Get(dbIdentity["id"], dbIdentity["kind"])
	if err != nil && err.Error() == "record not found" {
		directory.Add(dbIdentity["id"], dbIdentity["kind"], a.ToProto())
	} else {
		directory.Update(dbIdentity["id"], dbIdentity["kind"], a.ToProto())
	}
	return nil
}

func (a *Animal) Fill(msg proto.Message) error {
	// Fill the animal
	return a.FromProto(msg.(*xPb.Animal))
}

func (a *Animal) ToString() string {
	// Convert the animal to a string
	return fmt.Sprintf("Animal: %s\n", a.base.ToString())
}

func (a *Animal) ToStatus() map[string]interface{} {
	// Convert the animal to a status
	return a.base.ToStatus()
}

func (a *Animal) ToProto() *xPb.Animal {
	// Convert the animal to a proto
	msg := &xPb.Animal{}
	dbIdentity := a.base.GetDBIdentifier()
	id := dbIdentity["id"]
	msg.Id = &id
	kind := dbIdentity["kind"]
	msg.Kind = &kind
	name := a.base.Name
	msg.Name = &name
	age := int64(a.base.Age)
	msg.Age = &age
	deleted := a.base.Deleted
	msg.Deleted = &deleted
	description := a.base.Description
	msg.Description = &description
	cloned := a.base.Cloned
	msg.Cloned = &cloned
	clonedFromRef := a.base.ClonedFromRef.String()
	msg.ClonedFromRef = &clonedFromRef
	return msg
}

func (a *Animal) FromProto(msg *xPb.Animal) error {
	// Convert the animal from a proto
	id, err := uuid.Parse(msg.GetId())
	if err != nil {
		return err
	}
	a.base.ID = id
	a.base.Kind = msg.GetKind()
	a.base.Name = msg.GetName()
	a.base.Age = int(msg.GetAge())
	a.base.Deleted = msg.GetDeleted()
	a.base.Description = msg.GetDescription()
	a.base.Cloned = msg.GetCloned()
	clonedFromRef, err := uuid.Parse(msg.GetClonedFromRef())
	if err != nil {
		return err
	}
	a.base.ClonedFromRef = clonedFromRef
	return nil
}
