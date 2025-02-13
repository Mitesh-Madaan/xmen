package base

import (
	"fmt"

	"github.com/google/uuid"

	xMsg "lib/message"
)

type Base struct {
	Name          string          `json:"name"`
	ID            uuid.UUID       `json:"id"`
	Kind          string          `json:"kind"`
	Age           int             `json:"age,omitempty"`
	Deleted       bool            `json:"deleted,omitempty"`
	Description   string          `json:"description,omitempty"`
	MessageList   []*xMsg.Message `json:"message_list,omitempty"`
	Cloned        bool            `json:"cloned,omitempty"`
	ClonedFromRef uuid.UUID       `json:"cloned_from_ref,omitempty"`
}

func NewBase(name string, id uuid.UUID, kind string, age int, deleted bool, description string, message_list []*xMsg.Message, cloned bool, cloned_from_ref uuid.UUID) *Base {
	// Create a new base
	return &Base{
		Name:          name,
		ID:            id,
		Kind:          kind,
		Age:           age,
		Deleted:       deleted,
		Description:   description,
		MessageList:   message_list,
		Cloned:        cloned,
		ClonedFromRef: cloned_from_ref}
}

func (b *Base) ToString() string {
	// Convert the base to a string
	data := ""
	data += fmt.Sprintf("Name: %s ", b.Name)
	data += fmt.Sprintf("ID: %s ", b.ID.String())
	data += fmt.Sprintf("Kind: %s ", b.Kind)
	data += fmt.Sprintf("Age: %d ", b.Age)
	data += fmt.Sprintf("Deleted: %t ", b.Deleted)
	data += fmt.Sprintf("Description: %s ", b.Description)
	data += fmt.Sprintf("Cloned: %t ", b.Cloned)
	data += fmt.Sprintf("Cloned From Ref: %s ", b.ClonedFromRef.String())
	data += fmt.Sprintf("Message List: %v\n", xMsg.DumpMessages(b.MessageList))
	return data
}

func (b *Base) ToStatus() map[string]interface{} {
	// Convert the base to a status
	return map[string]interface{}{
		"id":   b.ID.String(),
		"kind": b.Kind,
		"data": map[string]interface{}{
			"name":            b.Name,
			"age":             b.Age,
			"deleted":         b.Deleted,
			"description":     b.Description,
			"cloned":          b.Cloned,
			"cloned_from_ref": b.ClonedFromRef.String(),
			"message_list":    xMsg.DumpMessages(b.MessageList),
		},
	}
}

func (b *Base) GetDBIdentifier() map[string]string {
	// Get the database identifier
	return map[string]string{
		"id":   b.ID.String(),
		"kind": b.Kind,
	}
}

func (b *Base) MarkDeleted() {
	// Mark the base as deleted
	b.Deleted = true
}

func (b *Base) Clone() *Base {
	// Clone the Base
	return NewBase(b.Name, uuid.New(), b.Kind, b.Age, false, b.Description, make([]*xMsg.Message, 0), true, b.ID)
}

func (b *Base) GetEditableFields() []string {
	// Get the editable fields
	return []string{"name", "description", "age"}
}

func (b *Base) PostEditables(editMap map[string]interface{}) {
	// Post the editable fields
	for key, value := range editMap {
		switch key {
		case "name":
			b.Name = value.(string)
		case "description":
			b.Description = value.(string)
		case "age":
			b.Age = value.(int)
		}
	}
}

func (b *Base) PreviewValidate() bool {
	// Validate the person
	messages := make([]*xMsg.Message, 0)
	if b.Name == "" {
		messages = append(messages, xMsg.NewMessage("name is required", "EmptyField"))
	}
	if b.ID == uuid.Nil {
		messages = append(messages, xMsg.NewMessage("id is required", "EmptyField"))
	}

	if len(messages) > 0 {
		b.MessageList = append(b.MessageList, messages...)
		return false
	}
	return true
}

func (b *Base) GetMessages() string {
	// Get the messages
	return xMsg.DumpMessages(b.MessageList)
}
