package config

import (
	"event-booking/common/constant"
	"event-booking/models"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database connection
func ConnectDB() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	DB = db
}

// Migrate function performs auto-migration for the database models
func Migrate() {
	// Run auto-migration
	err := DB.AutoMigrate(&models.User{}, &models.Event{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration completed successfully.")
}

func SeedUsers() {

	var count int64
	if DB.Model(&models.User{}).Count(&count); count > 0 {
		log.Println("Users already seeded!")
		return
	}

	// Hash the PIN
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to seed create hashed PIN!")
	}
	tempPassword := string(hashedPIN)

	users := []models.User{
		{Username: "HR1", Password: tempPassword, Role: constant.HR, FullName: "HR 1"},
		{Username: "HR2", Password: tempPassword, Role: constant.HR, FullName: "HR 2"},
		{Username: "Vendor1", Password: tempPassword, Role: constant.VENDOR, FullName: "Vendor 1"},
		{Username: "Vendor2", Password: tempPassword, Role: constant.VENDOR, FullName: "Vendor 2"},
	}

	for _, user := range users {
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v\n", user.Username, err)
		}
	}
	log.Println("Users table seeded successfully!")
}

func SeedEvents() {

	var count int64
	if DB.Model(&models.Event{}).Count(&count); count > 0 {
		log.Println("Events already seeded!")
		return
	}

	proposedDate := time.Now()
	proposedDates := fmt.Sprintf("%s,%s,%s", proposedDate.Format("2006-01-02"), proposedDate.AddDate(0, 0, 1).Format("2006-01-02"), proposedDate.AddDate(0, 0, 2).Format("2006-01-02"))
	location := "Jl. Kyai Maja No.43, Gunung, Kec. Kby. Baru, Kota Jakarta Selatan, Daerah Khusus Ibukota Jakarta 12120"
	eventName := "Vacine boost"
	events := []models.Event{
		{CompanyName: "ABC", ProposedDates: proposedDates, Location: location, EventName: eventName, Status: constant.PENDING, VendorID: 3, CreatedBy: 1, CreatedAt: time.Now()},
		{CompanyName: "DEF", ProposedDates: proposedDates, Location: location, EventName: eventName, Status: constant.PENDING, VendorID: 4, CreatedBy: 1, CreatedAt: time.Now()},
		{CompanyName: "GHI", ProposedDates: proposedDates, Location: location, EventName: eventName, Status: constant.PENDING, VendorID: 3, CreatedBy: 2, CreatedAt: time.Now()},
		{CompanyName: "JKL", ProposedDates: proposedDates, Location: location, EventName: eventName, Status: constant.PENDING, VendorID: 4, CreatedBy: 2, CreatedAt: time.Now()},
	}

	for _, event := range events {
		if err := DB.Create(&event).Error; err != nil {
			log.Printf("Failed to seed user %s: %v\n", event.CompanyName, err)
		}
	}
	log.Println("Events table seeded successfully!")
}
