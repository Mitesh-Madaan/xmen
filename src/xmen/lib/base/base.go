package base

import (
	xDb "lib/dbchef"
)

type Base interface {
	ToString() string
	ToStatus() map[string]interface{}
	PostEditableFields(objMap map[string]interface{}) error
	Clone() Base
	Create(dbSession *xDb.DBSession, objMap map[string]interface{}) error
	Update(dbSession *xDb.DBSession, objMap map[string]interface{}) error
	Delete(dbSession *xDb.DBSession) error
	Save(dbSession *xDb.DBSession, updates map[string]interface{}) error
}
