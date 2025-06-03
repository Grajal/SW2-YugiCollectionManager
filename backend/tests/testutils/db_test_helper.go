package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Grajal/SW2-YugiCollectionManager/backend/database"
	"github.com/Grajal/SW2-YugiCollectionManager/backend/models"

	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

// Create a new test DB, if it doesn't exist
func createTestDB() {

	defaultDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		"postgres",
		os.Getenv("PGPORT"),
	)

	db, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		log.Fatalf("Failed to connect to default database: %v", err)
	}
	defer db.Close()

	// Test DB name
	testDBName := os.Getenv("PGDATABASE")

	// Terminate all connections to the test database
	_, err = db.Exec(fmt.Sprintf(`
    -- Terminate all connections to the test database
    SELECT pg_terminate_backend(pg_stat_activity.pid)
    FROM pg_stat_activity
    WHERE pg_stat_activity.datname = '%s'
    AND pid <> pg_backend_pid();
   `, testDBName))

	if err != nil {
		log.Fatalf("Failed to terminate connections to the test database %s: %v", testDBName, err)
	}

	// Drop the test database if it exists
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName))
	if err != nil {
		log.Fatalf("Failed to drop test database %s: %v", testDBName, err)
	}

	// Create the test database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		log.Fatalf("Failed to create test database %s: %v", testDBName, err)
	} else {
		log.Printf("Database %s created successfully", testDBName)
	}
}

// TestDB holds the test database connection
var TestDB *gorm.DB

// SetupTestDatabase sets up the test database connection and runs migrations
func SetupTestDatabase() {

	createTestDB()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGPORT"),
	)

	var err error
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations (you can add your models here)
	if err := TestDB.AutoMigrate(
		models.User{},
		models.Card{},
		models.SpellTrapCard{},
		models.MonsterCard{},
		models.LinkMonsterCard{},
		models.PendulumMonsterCard{},
		models.UserCard{},
		models.Deck{},
		models.DeckCard{},
	); err != nil {
		log.Fatalf("Failed to auto migrate database schema: %v", err)
	}

}

// ResetTestDatabase resets the database after each test
func ResetTestDatabase() {
	TestDB.Exec("TRUNCATE TABLE items RESTART IDENTITY CASCADE")
}

// TearDownTestDatabase closes the database connection
func TearDownTestDatabase() {
	sqlDB, err := TestDB.DB()
	if err != nil {
		log.Fatalf("Failed to close the database: %v", err)
	}
	sqlDB.Close()
}

// PatchDatabase replaces the global database.DB with TestDB for testing
func PatchDatabase() {
	database.DB = TestDB
}

// UnpatchDatabase restores the original database.DB (if needed)
func UnpatchDatabase() {
	// Reset database.DB back to the development DB after testing
	database.DBConnect()
}

// InitializeTestSuite sets up the test environment and patches the database connection
func InitializeTestSuite(m *testing.M) {
	// Initialize the test database
	SetupTestDatabase()
	// Patch the global database.DB to use the TestDB
	PatchDatabase()
	// Run all tests
	code := m.Run()
	// Teardown the database connection
	TearDownTestDatabase()
	// Unpatch the database
	UnpatchDatabase()
	// Exit with the test run status code
	os.Exit(code)
}
