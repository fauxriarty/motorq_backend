package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	dsn := "user=postgres.wnugfqtxchhvoxtxnpha password=Adityachhed@1 host=aws-0-ap-south-1.pooler.supabase.com port=6543 dbname=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get sqlDB from GORM DB:", err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 5)

	DB = db

	fmt.Println("Database connection established")
}
