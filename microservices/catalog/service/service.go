package service

import (
	m "github.com/gosmo-devs/microsamples/catalog/model"
)

// ProductService manages product bussiness logic
type ProductService struct {
	dao m.ProductDAO
}

// GetProducts obtains products in batches
func (s *ProductService) GetProducts(offset int, limit int) (*m.ProductPage, *m.ProductError) {
	return s.dao.GetProducts(offset, limit)
}

// FindProduct finds a product by id
func (s *ProductService) FindProduct(id int) (*m.Product, *m.ProductError) {
	return s.dao.FindProduct(id)
}

// AddProduct adds a product
func (s *ProductService) AddProduct(p *m.Product) {
	s.dao.AddProduct(p)
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(id int) (bool, *m.ProductError) {
	return s.dao.DeleteProduct(id)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(id int, p *m.Product) (bool, *m.ProductError) {
	return s.dao.UpdateProduct(id, p)
}
