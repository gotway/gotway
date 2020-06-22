package model

import (
	"sort"
	"sync"
)

// StockList represents the stock of a list of products
type StockList struct {
	mStock      sync.Mutex
	mOutOfStock sync.Mutex
	mError      sync.Mutex
	Stock       []Stock      `json:"stock"`
	OutOfStock  []int        `json:"outOfStock"`
	Error       []StockError `json:"error"`
}

// AddStock adds a product with stock
func (sl *StockList) AddStock(s *Stock) {
	sl.mStock.Lock()
	defer sl.mStock.Unlock()
	sl.Stock = append(sl.Stock, *s)
}

// HandleError handles the error my updating the corresponding field
func (sl *StockList) HandleError(err *StockError) {
	if err.IsOutOfStockError() {
		sl.mOutOfStock.Lock()
		defer sl.mOutOfStock.Unlock()
		sl.OutOfStock = append(sl.OutOfStock, err.ProductID)
	}
	if err.IsInternalError() {
		sl.mError.Lock()
		defer sl.mError.Unlock()
		sl.Error = append(sl.Error, *err)
	}
}

// HasStock determines if there is stock available
func (sl *StockList) HasStock() bool {
	sl.mStock.Lock()
	defer sl.mStock.Unlock()
	return sl.Stock != nil && len(sl.Stock) > 0
}

// IsValid determines if a StockList is valid
func (sl *StockList) IsValid() bool {
	if !sl.HasStock() {
		return false
	}
	sl.mStock.Lock()
	defer sl.mStock.Unlock()
	for _, stock := range sl.Stock {
		if !stock.IsValid() {
			return false
		}
	}
	return true
}

// Get returns a copy of the current instance
func (sl *StockList) Get() StockList {
	sl.mStock.Lock()
	stock := make([]Stock, len(sl.Stock))
	copy(stock, sl.Stock)
	sort.Slice(stock, func(i, j int) bool {
		return stock[i].ProductID < stock[j].ProductID
	})
	sl.mStock.Unlock()

	sl.mOutOfStock.Lock()
	outOfStock := make([]int, len(sl.OutOfStock))
	copy(outOfStock, sl.OutOfStock)
	sort.Slice(outOfStock, func(i, j int) bool {
		return outOfStock[i] < outOfStock[j]
	})
	sl.mOutOfStock.Unlock()

	sl.mError.Lock()
	err := make([]StockError, len(sl.Error))
	copy(err, sl.Error)
	sort.Slice(err, func(i, j int) bool {
		return err[i].ProductID < err[j].ProductID
	})
	sl.mError.Unlock()

	return StockList{Stock: stock, OutOfStock: outOfStock, Error: err}
}
