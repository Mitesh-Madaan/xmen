// CRUD operations

package session

import (
	xError "gomike/error"
	xDb "lib/dbchef"
)

type Storable any

func CreateRecord[T Storable](dbSession *xDb.DBSession, obj T) error {
	err := dbSession.CreateRecord(&obj)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func ReadRecord[T Storable](dbSession *xDb.DBSession, objID string) (*T, error) {
	obj := new(T)
	err := dbSession.ReadRecord(map[string]interface{}{"id": objID}, obj)
	if err != nil {
		return nil, xError.NewDBError(err)
	}
	return obj, nil
}

func UpdateRecord[T Storable](dbSession *xDb.DBSession, obj T) error {
	err := dbSession.UpdateRecord(&obj)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}

func DeleteRecord[T Storable](dbSession *xDb.DBSession, obj T) error {
	err := dbSession.DeleteRecord(&obj)
	if err != nil {
		return xError.NewDBError(err)
	}
	return nil
}
