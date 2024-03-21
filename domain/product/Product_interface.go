package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"superindo/v1/model"
)

type Repository interface {
	Rep_begin_transaction() (*sqlx.Tx, error)
	Rep_commit_transaction(tx *sqlx.Tx) error
	Rep_rollback_transaction(tx *sqlx.Tx) error
	Rep_insert_product(tx *sqlx.Tx, payload model.Product) (*sqlx.Tx, error)
	Rep_get_product_list(offset, limit int) ([]model.Product, error)
	Rep_get_product_search_name_or_product_code(data string) (model.Product, error)
	Rep_get_product_filter_by_product_code(product_type string, offset, limit int) ([]model.Product, error)
	Rep_get_product_sort_by_key(key string, offset, limit int) ([]model.Product, error)
}

type Usecase interface {
	Use_post_product_register(c echo.Context, request model.Req_register_product) (interface{}, error)
	Use_get_product_list(c echo.Context) (interface{}, error)
	Use_get_product_search(c echo.Context) (interface{}, error)
	Use_get_product_filter(c echo.Context) (interface{}, error)
	Use_get_product_sort(c echo.Context) (interface{}, error)
}
