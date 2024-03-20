package model

import (
	"github.com/jmoiron/sqlx"
)

type Product struct {
	ID            int    `json:"-" db:"id"`
	Name          string `json:"name" db:"name"`
	Product_code  string `json:"product_code" db:"product_code"`
	Product_type  string `json:"product_type" db:"product_type"`
	Description   string `json:"description" db:"description"`
	Price         int    `json:"price" db:"price"`
	Register_date string `json:"register_date" db:"register_date"`
}

type Req_register_product struct {
	Name         string `json:"name" db:"name" validate:"required"`
	Product_type string `json:"product_type" db:"product_type" validate:"required"`
	Description  string `json:"description" db:"description" validate:"required"`
	Price        int    `json:"price" db:"price" validate:"required"`
}

type Query_proccessing struct {
	Tx  *sqlx.Tx
	Err error
}
