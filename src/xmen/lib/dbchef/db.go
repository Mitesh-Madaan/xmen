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
		err := s.conn.Migrator().CreateTable(model)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateRecords inserts multiple records into the database
func (s *DBSession) CreateRecords(model interface{}, records []interface{}) error {
	for _, record := range records {
		result := s.conn.Model(model).Create(record)
		if result.Error != nil {
			return result.Error
		}
		// Check if the record was created successfully
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
	}
	return nil
}

// ReadRecords retrieves records from the database based on the provided conditions
func (s *DBSession) ReadRecords(model interface{}, conditions map[string]interface{}, records interface{}) error {
	// Add notDeleted condition to avoid soft-deleted records
	conditions["Deleted"] = false
	result := s.conn.Model(model).Where(conditions).Find(records)
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
func (s *DBSession) UpdateRecords(model interface{}, conditions map[string]interface{}, updates map[string]interface{}) error {
	result := s.conn.Model(model).Where(conditions).Updates(updates)
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
func (s *DBSession) DeleteRecords(model interface{}, conditions map[string]interface{}) error {
	updates := map[string]interface{}{
		"Deleted": true,
	}
	return s.UpdateRecords(model, conditions, updates)
}

// Expose the session
var Session *DBSession

var connStr = "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"

func InitSession(connStr string) {
	Session = NewDBSession(connStr)
}

func GetSession() *DBSession {
	if Session == nil {
		InitSession(connStr)
	}
	return Session
}
