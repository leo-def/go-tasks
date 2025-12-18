package database

import (
    "log"
    "strings"

    "gorm.io/driver/postgres"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func ConnectDB(driver, dns string) Connection {
    var dialector gorm.Dialector

    switch strings.ToLower(driver) {
    case "postgres", "postgresql":
        dialector = postgres.Open(dns)
    case "sqlite", "sqlite3":
        dialector = sqlite.Open(dns)
    default:
        log.Fatalf("Unsupported database driver: %s. Supported drivers: postgres, sqlite", driver)
    }

    db, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        log.Fatal("Cannot connect to DB:", err)
    }
    // Ensure SQLite enforces foreign key constraints
    if strings.ToLower(driver) == "sqlite" || strings.ToLower(driver) == "sqlite3" {
        if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
            log.Printf("warning: failed to enable SQLite foreign_keys pragma: %v", err)
        }
    }
    return WrapConnection(db)
}
