package session

import (
	xDb "lib/dbchef"
)

var dbSession *xDb.DBSession

func GetDBSession(connStr string) *xDb.DBSession {
	if dbSession != nil {
		return dbSession
	}
	return InitDBSession(connStr)
}

func InitDBSession(connStr string) *xDb.DBSession {

	dbSession = xDb.GetSession(connStr)
	return dbSession
}
