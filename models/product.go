package models

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	Description   string         `gorm:"type:text" json:"description"`
	Price         float64        `gorm:"type:decimal(10,2);not null" json:"price" binding:"required"`
	StockQuantity int            `gorm:"not null;default:0" json:"stock_quantity" binding:"required"`
	Category      string         `gorm:"type:varchar(100)" json:"category"`
	ImageURL      string         `gorm:"type:text" json:"image_url"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
