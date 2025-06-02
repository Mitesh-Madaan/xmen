package base

type Base interface {
	ToString() string
	ToStatus() map[string]interface{}
	GetEditableFields() string
	Clone() Base
	Create(objMap map[string]interface{}) error
	Update(objMap map[string]interface{}) error
	Delete() error
	Save(updates map[string]interface{}) error
}
