package dbchef

import (
	"errors"
)

var db *Directory

type Directory struct {
	Records []Record
}

type Record struct {
	Id         string
	Kind       string
	ObjDetails []byte
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

func NewRecord(id, kind string, details []byte) Record {
	// Create a new record
	return Record{
		Id:         id,
		Kind:       kind,
		ObjDetails: details,
	}
}

// Get retrieves an object by its ID and kind from the Directory.
func (d *Directory) Get(id, kind string) ([]byte, error) {
	for _, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			return record.ObjDetails, nil
		}
	}
	return nil, errors.New("record not found")
}

// Add adds a new object to the Directory.
func (d *Directory) Add(id, kind string, details []byte) error {
	d.Records = append(d.Records, NewRecord(id, kind, details))
	return nil
}

// Update updates an existing object in the Directory.
func (d *Directory) Update(id, kind string, details []byte) error {
	for i, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			d.Records[i].ObjDetails = details
			return nil
		}
	}
	return errors.New("record not found")
}

// Delete removes an object from the Directory by its ID and kind.
func (d *Directory) Delete(id, kind string) error {
	for i, record := range d.Records {
		if record.Id == id && record.Kind == kind {
			d.Records = append(d.Records[:i], d.Records[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}
