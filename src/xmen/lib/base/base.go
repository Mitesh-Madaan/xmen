package base

type Base interface {
	ToString() string
	ToStatus() map[string]interface{}
	Clone() Base
	Create(objMap map[string]interface{}) error
	Update(objMap map[string]interface{}) error
	Delete() error
	Save() error
}
