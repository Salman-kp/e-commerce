package config

import (
	"fmt"
	"e-commerce/models"
)

// MigrateAll runs GORM auto migrations for all models
func MigrateAll() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.RefreshToken{},
		&models.Product{},
		&models.ProductProduction{},
		&models.CartItem{},
		&models.WishlistItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},
	)

	if err != nil {
		fmt.Println("❌ Migration failed:", err)
		return
	}
	fmt.Println("✅ All models migrated successfully!")
}
