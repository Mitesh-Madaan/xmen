package base

import (
	xDb "lib/dbchef"
)

type Base interface {
	ToString() string
	ToStatus() map[string]interface{}
	Clone() Base
	Create(dbSession *xDb.DBSession, objDetails []byte) error
	Update(dbSession *xDb.DBSession, objDetails []byte) error
	Delete(dbSession *xDb.DBSession) error
	Save(dbSession *xDb.DBSession, updates map[string]interface{}) error
}
