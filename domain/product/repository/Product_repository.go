package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"superindo/v1/domain/product"
	"superindo/v1/model"
)

type Product_repository struct {
	sqlx *sqlx.DB
}

func NewProduct_repository(sqlx *sqlx.DB) product.Repository {
	return &Product_repository{sqlx}
}

func (IN *Product_repository) Rep_begin_transaction() (*sqlx.Tx, error) {
	tx, err := IN.sqlx.Beginx()
	return tx, err
}

func (IN *Product_repository) Rep_commit_transaction(tx *sqlx.Tx) error {

	err := tx.Commit()
	if err != nil {
		return IN.Rep_rollback_transaction(tx)
	}
	return nil
}

func (IN *Product_repository) Rep_rollback_transaction(tx *sqlx.Tx) error {
	var err error
	if err = tx.Rollback(); err != nil {
		return err
	}
	return err
}

func (IN *Product_repository) Rep_insert_product(tx *sqlx.Tx, payload model.Product) (*sqlx.Tx, error) {
	var err error
	query := `
		INSERT INTO product (name, product_code, product_type, description, price)
		VALUES ($1, $2, $3, $4,$5) ;
	`
	if _, err = tx.Exec(query,
		payload.Name,
		payload.Product_code,
		payload.Product_type,
		payload.Description,
		payload.Price,
	); err != nil {
		return tx, err
	}

	return tx, nil
}

func (IN *Product_repository) Rep_get_product_list(offset, limit int) ([]model.Product, error) {
	var products []model.Product
	var err error
	query := `
		SELECT   name, product_type, description, price, to_char(register_date, 'YYYY-MM-DD') as register_date
		FROM product
		ORDER BY id
		OFFSET $1 LIMIT $2 ;
	`

	if err = IN.sqlx.Select(&products, query, offset, limit); err != nil {
		return nil, err
	}

	return products, nil
}

func (IN *Product_repository) Rep_get_product_search_name_or_product_code(query string) (model.Product, error) {
	var err error
	var data model.Product

	// SQL untuk mencari data product by name atau product_code
	sql := `
		SELECT  name, product_code, product_type, description, price,  to_char(register_date, 'YYYY-MM-DD') as register_date 
		FROM product
		WHERE name ILIKE $1 OR product_code = $2 ;
	`

	if err = IN.sqlx.Get(&data, sql, "%"+query+"%", query); err != nil {
		return data, err
	}

	return data, nil
}

func (IN *Product_repository) Rep_get_product_filter_by_product_code(product_type string, offset, limit int) ([]model.Product, error) {
	var err error
	var data []model.Product

	sql := `
		SELECT  name, product_code, product_type, description, price, to_char(register_date, 'YYYY-MM-DD') as register_date  
		FROM product
		WHERE product_type = $1 OFFSET $2 LIMIT $3 ;
	`

	if err = IN.sqlx.Select(&data, sql, product_type, offset, limit); err != nil {
		return data, err
	}

	return data, nil
}

func (IN *Product_repository) Rep_get_product_sort_by_key(key string, offset, limit int) ([]model.Product, error) {
	var err error
	var data []model.Product

	sql := fmt.Sprintf(`
		SELECT name, product_code, product_type, description, price, to_char(register_date, 'YYYY-MM-DD') as register_date
		FROM product
		ORDER BY %s DESC
		OFFSET $1 LIMIT $2`, key)

	if err = IN.sqlx.Select(&data, sql, offset, limit); err != nil {
		return data, err
	}

	return data, nil
}
