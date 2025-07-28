package dbchef

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

type DBSession struct {
	conn *gorm.DB
}

// NewDBSession initializes a singleton DB connection
func NewDBSession(connStr string) *DBSession {
	once.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
	})
	return &DBSession{conn: db}
}

// SeedTables seeds the database with initial data for the provided models
func (s *DBSession) SeedTables(models []interface{}) error {
	for _, model := range models {
		if !s.conn.Migrator().HasTable(model) {
			err := s.conn.Migrator().CreateTable(model)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateRecords inserts multiple records into the database
func (s *DBSession) CreateRecord(record interface{}) error {
	result := s.conn.Model(record).Create(record)
	if result.Error != nil {
		return result.Error
	}
	// Check if the record was created successfully
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ReadRecords retrieves records from the database based on the provided conditions
func (s *DBSession) ReadRecord(conditions map[string]interface{}, record interface{}) error {
	result := s.conn.Model(record).Find(record, conditions)
	if result.Error != nil {
		return result.Error
	}
	// Check if any records were found
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateRecords updates records in the database based on the provided conditions
func (s *DBSession) UpdateRecord(record interface{}) error {
	result := s.conn.Model(record).Updates(record)
	if result.Error != nil {
		return result.Error
	}
	// Check if any records were updated
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteRecords deletes records from the database based on the provided conditions
func (s *DBSession) DeleteRecord(record interface{}) error {
	result := s.conn.Model(record).Delete(record)
	if result.Error != nil {
		return result.Error
	}
	// Check if any records were deleted
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
