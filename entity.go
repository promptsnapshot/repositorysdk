package repository

import (
	"github.com/google/uuid"
	utils "github.com/thinc-org/newbie-utils"
	"gorm.io/gorm"
	"time"
)

type Base struct {
	ID        *uuid.UUID     `json:"id" gorm:"primary_key"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp;autoUpdateTime:nano"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamp"`
}

func (b *Base) BeforeCreate(_ *gorm.DB) error {
	if b.ID == nil {
		b.ID = utils.UUIDAdr(uuid.New())
	}

	return nil
}

type PaginationMetadata struct {
	ItemsPerPage int
	ItemCount    int
	TotalItem    int
	CurrentPage  int
	TotalPage    int
}

func (p *PaginationMetadata) GetOffset() int {
	return (p.GetCurrentPage() - 1) * p.GetItemPerPage()
}

func (p *PaginationMetadata) GetItemPerPage() int {
	if p.ItemsPerPage < 10 {
		p.ItemsPerPage = 10
	}
	if p.ItemsPerPage > 100 {
		p.ItemsPerPage = 100
	}

	return p.ItemsPerPage
}

func (p *PaginationMetadata) GetCurrentPage() int {
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
	return p.CurrentPage
}
