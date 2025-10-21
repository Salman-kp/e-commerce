package models

import "time"

type WishlistItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_product,unique" json:"-"`
	ProductID uint      `gorm:"not null;index:idx_user_product,unique" json:"-"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}