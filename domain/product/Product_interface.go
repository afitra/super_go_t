package product

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Rep_begin_transaction() (*sqlx.Tx, error)
	Rep_commit_transaction(tx *sqlx.Tx) error
	Rep_rollback_transaction(tx *sqlx.Tx) error
}

type Usecase interface {
	//Use_post_product(c echo.Context) (interface{}, error)
}
