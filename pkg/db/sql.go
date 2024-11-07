package db

import (
	"fmt"
	"spam-search/pkg/config"
	"spam-search/pkg/users"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type sqlConnection struct {
	connection *gorm.DB
}

func NewSQLDB() (*sqlConnection, error) {
	dbConfig := config.DBconfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true", dbConfig.Username, dbConfig.Password, dbConfig.URL, dbConfig.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	return &sqlConnection{connection: db}, nil
}

func (s *sqlConnection) GetDatabse() *gorm.DB {
	return s.connection
}

func (s *sqlConnection) Migrate(development bool) {
	s.connection.AutoMigrate(&users.User{})
}
