package dbchef

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

var db *Directory

type Directory struct {
	Records []Record
}

type Entity[Proto proto.Message] interface {
	ToProto() Proto
	FromProto(p Proto) error
}

type Record struct {
	Id   string
	Kind string
	Msg  proto.Message
}

func GetDirectory() *Directory {
	// Get the directory
	if db != nil {
		return db
	}
	db = createDirectory()
	return db
}

func createDirectory() *Directory {
	// Create a new directory
	directory := Directory{}
	directory.Records = make([]Record, 0)
	return &directory
}

func NewRecord(id, kind string, msg proto.Message) Record {
	// Create a new record
	return Record{
		Id:   id,
		Kind: kind,
		Msg:  msg,
	}
}

// Get retrieves a proto.Message by its ID and kind from the Directory.
func (d *Directory) Get(id, kind string) (proto.Message, error) {
	for _, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			return record.Msg, nil
		}
	}
	return nil, errors.New("record not found")
}

// Add adds a new proto.Message to the Directory.
func (d *Directory) Add(id, kind string, msg proto.Message) error {
	d.Records = append(d.Records, NewRecord(id, kind, msg))
	return nil
}

// Update updates an existing proto.Message in the Directory.
func (d *Directory) Update(id, kind string, msg proto.Message) error {
	for i, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			d.Records[i] = NewRecord(id, kind, msg)
			return nil
		}
	}
	return errors.New("record not found")
}

// Delete removes a proto.Message from the Directory by its ID and kind.
func (d *Directory) Delete(id, kind string) error {
	for i, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			d.Records = append(d.Records[:i], d.Records[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}
