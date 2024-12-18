package models

import (
	"time"
)

type Event struct {
	ID            uint `gorm:"primaryKey"`
	CompanyName   string
	ProposedDates string // Comma-separated
	Location      string
	EventName     string
	Status        string // Pending, Approved, Rejected
	Remarks       string
	ConfirmedDate string
	VendorID      uint
	CreatedBy     uint
	CreatedAt     time.Time
}

type EventWithVendorName struct {
	ID            uint `gorm:"primaryKey"`
	CompanyName   string
	ProposedDates string // Comma-separated
	Location      string
	EventName     string
	Status        string // Pending, Approved, Rejected
	Remarks       string
	ConfirmedDate string
	VendorID      uint
	CreatedBy     uint
	CreatedAt     time.Time
	VendorName    string
}
