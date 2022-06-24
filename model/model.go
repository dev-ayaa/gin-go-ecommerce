package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName     *string            `json: "first_name" validate: "required,min=2,max=30"`
	LastName       *string            `json: "last_name" validate:"required,min=2,max=30"`
	Password       *string            `json: "password" validate:"required,min=8"`
	Email          *string            `json: "email" validate: "email,required"`
	Phone          *string            `json: "phone"`
	Token          *string            `json: "token"`
	RefreshToken   *string            `json:"refresh_token"`
	UserID         string             `json: "user_id"`
	UserCart       []ProductUser      `json: "user_cart" bson:"user_cart"`
	AddressDetails []Address          `json: "address_details" bson:"address_details"`
	OrderStatus    []Order            `json: "order_status" bson :"order_status"`
	CreatedAt      time.Time          `json: "created_at"`
	UpdatedAt      time.Time          `json: "updated_at"`
}
type Product struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json:"product_name"`
	Price       *uint64            `json: "price"`
	Rating      uint               `json: "rating"`
	Image       *string            `json: "image"`
}

type ProductUser struct {
	ProductID   primitive.ObjectID `bson:"_id"`
	ProductName *string            `json: "product_name" bson:"product_name"`
	Price       *int               `json: "price" bson:"price"`
	Rating      *uint              `json: "rating" bson:"rating"`
	Image       *string            `json: "image" bson:"image"`
}

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id"`
	Order_Cart    []ProductUser      `json: "order_cart" bson:"order_cart"`
	Price         int                `json: "price" bson:"price"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
	Discount      int                `json: "discount" bson:"discount"`
	OrderedAt     time.Time          `json: "ordered_at" bson:"ordered_at"`
}

type Payment struct {
	Digital bool `json:"digital"`
	COD     bool `json:"cod"`
}

type Address struct {
	AddressID primitive.ObjectID `bson:"_id"`
	PinCode   *string            `json: "pin_code" bson:"pin_code"`
	City      *string            `json: "city" bson:"city"`
	Street    *string            `json: "street" bson:"street"`
	House     *string            `json: "house" bson:"house"`
}
