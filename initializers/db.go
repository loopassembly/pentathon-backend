package initializers

import (
	// "fmt"
	"log"
	"os"

	"github.com/loopassembly/pentathon-backend/models"
	// "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
// UNIQUE constraint failed: emails.email

var DB *gorm.DB

func ConnectDB(config *Config) {
    var err error

    // Set the SQLite3 database path
    // dbPath := config.DBPath
    // DB
    // dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	// dsn := fmt.Sprintf("host=postgres port=6500 user=admin password=password123 dbname=golang_fiber sslmode=disable TimeZone=Asia/Kolkata")


	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{}) //?sqlite
// 
if err != nil {
	log.Fatal("Failed to connect to the Database! \n", err.Error())
	os.Exit(1)
}

log.Println("Running Migrations")// status
err = DB.AutoMigrate(&models.Email{})
if err != nil {
	log.Fatal("Migration Failed:  \n", err.Error())
	os.Exit(1)
}

    log.Println("🚀 Connected Successfully to the Database")
}