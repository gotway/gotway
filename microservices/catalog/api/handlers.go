package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	m "github.com/gotway/gotway/microservices/catalog/model"
	s "github.com/gotway/gotway/microservices/catalog/service"
)

var productService s.ProductService

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	setHeaders := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache-Tags", "ecommerce")
	}
	q := r.URL.Query()
	offsetStr := q.Get("offset")
	limitStr := q.Get("limit")
	offset, limit, err := processPaginationParams(offsetStr, limitStr)

	if err != nil {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productPage, productErr := productService.GetProducts(offset, limit)
	if productErr != nil {
		setHeaders(w)
		handleError(w, productErr)
		return
	}

	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productPage)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product m.Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	if !product.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productService.AddProduct(&product)
	res := m.ProductCreated{ID: product.ID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	setHeaders := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "s-maxage=10")
	}
	id, err := getIDparam(r)

	if err != nil {
		setHeaders(w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, productErr := productService.FindProduct(id)
	if productErr != nil {
		setHeaders(w)
		handleError(w, productErr)
		return
	}

	setHeaders(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIDparam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, productErr := productService.DeleteProduct(id)
	if productErr != nil {
		handleError(w, productErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := getIDparam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product m.Product
	_ = json.NewDecoder(r.Body).Decode(&product)
	if !product.IsValid() {
		w.WriteHeader(http.StatusBadGateway)
	}
	_, productErr := productService.UpdateProduct(id, &product)
	if err != nil {
		handleError(w, productErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIDparam(r *http.Request) (int, error) {
	params := mux.Vars(r)
	return strconv.Atoi(params["id"])
}

func handleError(w http.ResponseWriter, err *m.ProductError) {
	log.Print(err)
	w.WriteHeader(err.Code)
}

func processPaginationParams(offsetStr string, limitStr string) (int, int, error) {
	offset, err := processIntParam(offsetStr, 0)
	if err != nil {
		return 0, 0, err
	}
	limit, err := processIntParam(limitStr, 10)
	if err != nil {
		return 0, 0, err
	}
	if offset > limit {
		return 0, 0, errors.New("Offset cannot be greater than limit")
	}
	return offset, limit, nil
}

func processIntParam(paramStr string, defaultValue int) (int, error) {
	if len(paramStr) == 0 {
		return defaultValue, nil
	}
	param, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, err
	}
	if param < 0 {
		return 0, errors.New("Param cannot not be negative")
	}
	return param, nil
}
