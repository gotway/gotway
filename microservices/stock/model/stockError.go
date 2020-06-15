package model

import (
	"errors"
	"fmt"
	"net/http"
)

// StockError model
type StockError struct {
	Code      int
	ProductID int
	Err       error
}

func (err *StockError) Error() string {
	return fmt.Sprintf("Code: %d | ProductID: %d | Error: %s", err.Code, err.ProductID, err.Err)
}

// IsOutOfStockError checks if an error is an out of stock error
func (err *StockError) IsOutOfStockError() bool {
	return err.Code == outOfStockErrorCode
}

// IsInternalError checks if an error is an internal error
func (err *StockError) IsInternalError() bool {
	return err.Code == internalErrorCode
}

var outOfStockErrorCode = http.StatusNotFound

// OutOfStockError is raised when a product has no stock
func OutOfStockError(productID int) *StockError {
	return &StockError{Code: outOfStockErrorCode, ProductID: productID, Err: errors.New("Out of stock")}
}

var internalErrorCode = http.StatusInternalServerError

// InternalError returns an internal error object
func InternalError(productID int, err error) *StockError {
	return &StockError{Code: http.StatusInternalServerError, ProductID: productID, Err: err}
}
