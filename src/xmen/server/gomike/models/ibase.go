package models

import (
	"fmt"
	xBase "lib/base"
	xDb "lib/dbchef"
)

type IBase interface {
	ToString() string
	PreviewValidate() bool
	Save() error
	Fill([]byte) error
	Edit(map[string]interface{}) error
	ToStatus() map[string]interface{}
	Delete() error
	Clone() IBase
	GetBase() *xBase.Base
}

type BaseConstructor func(*xBase.Base) IBase

var BaseMapper = map[string]BaseConstructor{
	"Person": NewPerson,
	"Animal": NewAnimal,
}

func GetBase(b *xBase.Base) IBase {
	// Get the base
	return BaseMapper[b.Kind](b)
}

func GetObjectFromDB(id, kind string) (IBase, error) {
	// Get the object from the DB
	directory := xDb.GetDirectory()
	if directory == nil {
		return nil, fmt.Errorf("Directory not found")
	}

	details, err := directory.Get(id, kind)
	if err != nil {
		return nil, err
	}

	b := &xBase.Base{
		Kind: kind,
	}
	obj := GetBase(b)
	err = obj.Fill(details)
	if err != nil {
		fmt.Printf("Error filling base from details: %v\n", err)
		return nil, err
	}
	return obj, nil
}
