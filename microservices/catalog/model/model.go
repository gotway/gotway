package model

import (
	"fmt"
	"net/http"
)

// Product model
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Color string `json:"color"`
	Size  string `json:"size"`
}

// IsValid checks whether a product is valid
func (p *Product) IsValid() bool {
	return len(p.Name) > 0 && p.Price > 0 && len(p.Color) > 0 && len(p.Size) > 0
}

// ProductError model
type ProductError struct {
	Code    int
	Message string
}

func (err *ProductError) Error() string {
	return fmt.Sprintf("Code: %d | Message: %s", err.Code, err.Message)
}

// ProductCreated model
type ProductCreated struct {
	ID int `json:"id"`
}

var notFoundError = &ProductError{Code: http.StatusNotFound, Message: "Not found"}

// ProductPage model
type ProductPage struct {
	Products   []Product `json:"products"`
	TotalCount int       `json:"totatCount"`
}
