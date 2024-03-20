package usecase

import "superindo/v1/domain/product"

type Product_usecase struct {
	source             string
	product_repository product.Repository
}

func NewProduct_usecase(product_repository product.Repository) product.Usecase {
	return &Product_usecase{"Product_usecase", product_repository}
}
