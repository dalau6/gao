package clients

import (
	"clients/gen/client"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Client *client.ClientManagement

// InitDB is the function that starts a database file and table structures
// if not created then returns db object for next functions
func InitDB() *gorm.DB {
	// Opening file
	db, err := gorm.Open("sqlite3", "./data.db")
	// Display SQL queries
	db.LogMode(true)

	// Error
	if err != nil {
		panic(err)
	}
	// Creating the table if it doesn't exist
	var TableStruct = client.ClientManagement{}
	if !db.HasTable(TableStruct) {
		db.CreateTable(TableStruct)
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(TableStruct)
	}

	return db
}

// GetClient retrieves one client by its ID
func GetClient(clientID string) (client.ClientManagement, error) {
	db := InitDB()
	defer db.Close()

	var clients client.ClientManagement
	db.Where("client_id = ?", clientID).First(&clients)
	return clients, err
}

// CreateClient created a client row in DB
func CreateClient(client Client) error {
	db := InitDB()
	defer db.Close()
	err := db.Create(&client).Error
	return err
}

// ListClients retrieves the clients stored in Database
func ListClients() (client.ClientManagementCollection, error) {
	db := InitDB()
	defer db.Close()
	var clients client.ClientManagementCollection
	err := db.Find(&clients).Error
	return clients, err
}
