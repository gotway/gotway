package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	m "github.com/gosmo-devs/microsamples/stock/model"
	s "github.com/gosmo-devs/microsamples/stock/service"

	"github.com/gorilla/mux"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func upsertStock(w http.ResponseWriter, r *http.Request) {
	productID, err := getProductID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var stockData m.StockData
	_ = json.NewDecoder(r.Body).Decode(&stockData)
	if !stockData.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stock := stockData.ToStock(productID)
	resultStock, stockError := s.UpsertStock(productID, &stock)
	if stockError != nil {
		handleError(w, stockError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultStock)
}

func upsertStockList(w http.ResponseWriter, r *http.Request) {
	var stockList m.StockList
	_ = json.NewDecoder(r.Body).Decode(&stockList)
	if !stockList.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	resultStockList := s.UpsertStockList(stockList.Stock)
	if !resultStockList.HasStock() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&resultStockList)
}

func getStock(w http.ResponseWriter, r *http.Request) {
	productID, err := getProductID(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stock, stockErr := s.GetStock(productID)
	if stockErr != nil {
		handleError(w, stockErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stock)
}

func getStockList(w http.ResponseWriter, r *http.Request) {
	productIDs := getProductIDs(r)
	if productIDs == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stockList := s.GetStockList(productIDs)
	if !stockList.HasStock() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&stockList)
}

func getProductID(r *http.Request) (int, error) {
	params := mux.Vars(r)
	productID, err := strconv.Atoi(params["id"])
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func getProductIDs(r *http.Request) []int {
	q := r.URL.Query()
	requestProductIDs := q["productId"]
	if requestProductIDs == nil {
		return nil
	}
	productIDs := []int{}
	for _, val := range requestProductIDs {
		if val == "" {
			continue
		}
		productID, err := strconv.Atoi(val)
		if err == nil {
			productIDs = append(productIDs, productID)
		}
	}
	if len(productIDs) == 0 {
		return nil
	}
	return productIDs
}

func handleError(w http.ResponseWriter, err *m.StockError) {
	log.Print(err)
	w.WriteHeader(err.Code)
}
