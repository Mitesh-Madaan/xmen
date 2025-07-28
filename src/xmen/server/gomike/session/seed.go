package session

import (
	"fmt"

	xError "gomike/error"
	xModels "gomike/models"
	xDb "lib/dbchef"
)

// SeedTables creates and seeds Person and Animal tables
func SeedTables(dbSession *xDb.DBSession) error {
	models := []interface{}{
		&xModels.Person{},
		&xModels.Animal{},
	}

	err := dbSession.SeedTables(models)
	if err != nil {
		err := fmt.Errorf("failed to seed tables: %w", err)
		return xError.NewDBError(err)
	}
	return nil
}
