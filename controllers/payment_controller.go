package controllers

import (
	"e-commerce/config"
	"e-commerce/models"
	"e-commerce/services"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

// PaymentResponse DTO
type PaymentResponse struct {
	ID        uint      `json:"id"`
	PaymentID string    `json:"payment_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Gateway   string    `json:"gateway"`
	OrderID   uint      `json:"order_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// POST /payments/create - Create Stripe PaymentIntent
func CreatePaymentIntent(c *gin.Context) {
	type RequestBody struct {
		OrderID uint `json:"order_id"`
	}
	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var order models.Order
	if err := config.DB.Preload("OrderItems.Product").Preload("User").First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already initiated or order not pending"})
		return
	}

	var existingPayment models.Payment
	if err := config.DB.First(&existingPayment, "order_id = ? AND status = ?", order.ID, "pending").Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pending payment already exists"})
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.TotalAmount * 100)),
		Currency: stripe.String(string(stripe.CurrencyINR)),
	}
	params.AddMetadata("order_id", strconv.Itoa(int(order.ID)))

	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payment := models.Payment{
		OrderID:   order.ID,
		Gateway:   "Stripe",
		PaymentID: pi.ID,
		Amount:    order.TotalAmount,
		Status:    "pending",
	}
	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	paymentResp := PaymentResponse{
		ID:        payment.ID,
		PaymentID: payment.PaymentID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		Gateway:   payment.Gateway,
		OrderID:   payment.OrderID,
		CreatedAt: payment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"payment":       paymentResp,
		"client_secret": pi.ClientSecret,
	})
}

// PUT /payments/:payment_id/update - Update payment status
func UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("payment_id")
	type Body struct {
		Status string `json:"status"`
	}
	var body Body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var payment models.Payment
	if err := config.DB.Preload("Order.OrderItems.Product").Preload("Order.User").
		First(&payment, "payment_id = ?", paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment already processed"})
		return
	}
	if body.Status != "succeeded" && body.Status != "failed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	payment.Status = body.Status
	if err := config.DB.Save(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
		return
	}

	switch body.Status {
	case "succeeded":
		for _, item := range payment.Order.OrderItems {
			if item.Product.StockQuantity >= item.Quantity {
				item.Product.StockQuantity -= item.Quantity
				config.DB.Save(&item.Product)
			}
		}
		config.DB.Where("user_id = ?", payment.Order.UserID).Delete(&models.CartItem{})
		payment.Order.Status = "processing"
		config.DB.Save(&payment.Order)
	case "failed":
		payment.Order.Status = "failed"
		config.DB.Save(&payment.Order)
	}

	orderItemsResp := []services.OrderItemResponse{}
	for _, item := range payment.Order.OrderItems {
		orderItemsResp = append(orderItemsResp, services.OrderItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	orderResp := services.OrderResponse{
		ID:          payment.Order.ID,
		TotalAmount: payment.Order.TotalAmount,
		Address:     payment.Order.Address,
		Status:      payment.Order.Status,
		CreatedAt:   payment.Order.CreatedAt,
		UserName:    payment.Order.User.FullName,
		Items:       orderItemsResp,
	}

	paymentResp := PaymentResponse{
		ID:        payment.ID,
		PaymentID: payment.PaymentID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		Gateway:   payment.Gateway,
		OrderID:   payment.OrderID,
		CreatedAt: payment.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: payment.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, gin.H{
		"payment": paymentResp,
		"order":   orderResp,
	})
}
