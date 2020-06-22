package service

import (
	"strconv"
	"sync"

	m "github.com/gosmo-devs/microsamples/stock/model"
	"github.com/gosmo-devs/microsamples/stock/redis"
)

// UpsertStock upserts the stock of a product
func UpsertStock(productID int, stock *m.Stock) (*m.Stock, *m.StockError) {
	key := getKey(productID)
	ttl, err := redis.Set(key, stock.Units, stock.TTL)
	if err != nil {
		return nil, m.InternalError(productID, err)
	}
	resultStock := &m.Stock{ProductID: productID, Units: stock.Units, TTL: ttl}
	return resultStock, nil
}

// UpsertStockList upserts the stock of a list of products
func UpsertStockList(stockList []m.Stock) m.StockList {
	var wg sync.WaitGroup
	var sl m.StockList

	wg.Add(len(stockList))
	for _, stock := range stockList {
		go func(s m.Stock) {
			defer wg.Done()
			resultStock, err := UpsertStock(s.ProductID, &s)
			if err == nil {
				sl.AddStock(resultStock)
			} else {
				sl.HandleError(err)
			}
		}(stock)
	}
	wg.Wait()

	return sl.Get()
}

// GetStock gets the stock of a product
func GetStock(productID int) (*m.Stock, *m.StockError) {
	unitsChan := make(chan stockResult)
	ttlChan := make(chan stockResult)

	go getStockUnits(productID, unitsChan)
	go getTTL(productID, ttlChan)

	units, ttl := <-unitsChan, <-ttlChan
	if units.err != nil {
		return nil, units.err
	}
	if ttl.err != nil {
		return nil, units.err
	}

	stock := &m.Stock{ProductID: productID, Units: units.val, TTL: ttl.val}
	return stock, nil
}

// GetStockList gets the stock of a list of products
func GetStockList(productIDs []int) m.StockList {
	var wg sync.WaitGroup
	var sl m.StockList

	wg.Add(len(productIDs))
	for _, productID := range productIDs {
		go func(id int) {
			defer wg.Done()
			stock, err := GetStock(id)
			if err == nil {
				sl.AddStock(stock)
			} else {
				sl.HandleError(err)
			}
		}(productID)
	}
	wg.Wait()

	return sl.Get()
}

func getKey(productID int) string {
	return strconv.Itoa(productID)
}

type stockResult struct {
	val int
	err *m.StockError
}

func getStockUnits(productID int, c chan stockResult) {
	key := getKey(productID)
	defer close(c)
	val, err := redis.Get(key)
	if err != nil {
		c <- stockResult{0, m.OutOfStockError(productID)}
		return
	}
	units, parseErr := strconv.Atoi(val)
	if parseErr != nil {
		c <- stockResult{0, m.InternalError(productID, parseErr)}
		return
	}
	if units <= 0 {
		c <- stockResult{0, m.OutOfStockError(productID)}
		return
	}
	c <- stockResult{units, nil}
}

func getTTL(productID int, c chan stockResult) {
	key := getKey(productID)
	defer close(c)
	ttl, err := redis.TTL(key)
	if err != nil {
		c <- stockResult{0, m.OutOfStockError(productID)}
		return
	}
	c <- stockResult{ttl, nil}
}
