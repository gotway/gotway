package model

// StockData represents the data relative to stock
type StockData struct {
	Units int `json:"units"`
	TTL   int `json:"ttl"`
}

// IsValid validates a StockData
func (sd *StockData) IsValid() bool {
	return sd.Units > 0 && sd.TTL > 0
}

// ToStock transforms a StockData into a Stock object
func (sd *StockData) ToStock(productID int) Stock {
	return Stock{productID, sd.Units, sd.TTL}
}

// Stock represents the stock units of a product
type Stock struct {
	ProductID int `json:"productId"`
	Units     int `json:"units"`
	TTL       int `json:"ttl"`
}

// IsValid determines if a Stock is valid
func (s *Stock) IsValid() bool {
	return s.ProductID > 0 && s.Units > 0 && s.TTL > 0
}
