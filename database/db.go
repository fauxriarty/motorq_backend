package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	dsn := "user=postgres.wnugfqtxchhvoxtxnpha password=Adityachhed@1 host=aws-0-ap-south-1.pooler.supabase.com port=6543 dbname=postgres"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	fmt.Println("Database connection established")
}
