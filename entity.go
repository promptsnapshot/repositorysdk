package repositorysdk

import (
	"github.com/google/uuid"
	gosdk "github.com/thinc-org/newbie-gosdk"
	"gorm.io/gorm"
	"time"
)

// Base is a struct that holds common fields for database tables, including the ID, creation and update timestamps,
// and soft deletion timestamp.
type Base struct {
	ID        *uuid.UUID     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamp"`
}

// BeforeCreate is a GORM callback that generates a new UUID for the ID field before creating a new record in the database.
func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == nil {
		b.ID = gosdk.UUIDAdr(uuid.New())
	}

	return nil
}

// BaseHardDelete is a struct that holds common fields for database tables, including the ID and creation and update timestamps,
// but excludes soft deletion timestamp.
type BaseHardDelete struct {
	ID        *uuid.UUID `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;autoCreateTime:nano"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp;autoUpdateTime:nano"`
}

// BeforeCreate is a GORM callback that generates a new UUID for the ID field before creating a new record in the database.
func (b *BaseHardDelete) BeforeCreate(_ *gorm.DB) error {
	if b.ID == nil {
		b.ID = gosdk.UUIDAdr(uuid.New())
	}

	return nil
}

// PaginationMetadata is a struct that holds pagination metadata including the number of items per page, the current page,
// the total number of items, and the total number of pages.
type PaginationMetadata struct {
	ItemsPerPage int
	ItemCount    int
	TotalItem    int
	CurrentPage  int
	TotalPage    int
}

// GetOffset is a method that calculates the offset for the current page based on the number of items per page.
func (p *PaginationMetadata) GetOffset() int {
	return (p.GetCurrentPage() - 1) * p.GetItemPerPage()
}

// GetItemPerPage is a method that returns the number of items per page, ensuring that the value is within a certain range.
func (p *PaginationMetadata) GetItemPerPage() int {
	if p.ItemsPerPage < 10 {
		p.ItemsPerPage = 10
	}
	if p.ItemsPerPage > 100 {
		p.ItemsPerPage = 100
	}

	return p.ItemsPerPage
}

// GetCurrentPage is a method that returns the current page, ensuring that the value is at least 1.
func (p *PaginationMetadata) GetCurrentPage() int {
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
	return p.CurrentPage
}

// CalItemPerPage is a method that calculates the number of items per page based on the total number of items and the
// desired number of items per page.
func (p *PaginationMetadata) CalItemPerPage() {
	if p.ItemCount < p.ItemsPerPage {
		p.ItemsPerPage = p.ItemCount
	}
}
