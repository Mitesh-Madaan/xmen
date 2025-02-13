package models

import (
	"fmt"
	xBase "lib/base"
	xDb "lib/dbchef"

	"google.golang.org/protobuf/proto"
)

type IBase interface {
	ToString() string
	PreviewValidate() bool
	Save() error
	Fill(msg proto.Message) error
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

func GetBaseFromProto(msg proto.Message, kind string) IBase {
	// Get the base from proto
	b := &xBase.Base{
		Kind: kind,
	}
	obj := GetBase(b)
	obj.Fill(msg)
	return obj
}

func GetObjectFromDB(id, kind string) (IBase, error) {
	// Get the object from the DB
	directory := xDb.GetDirectory()
	if directory == nil {
		return nil, fmt.Errorf("Directory not found")
	}

	msg, err := directory.Get(id, kind)
	if err != nil {
		return nil, err
	}

	return GetBaseFromProto(msg, kind), nil
}
