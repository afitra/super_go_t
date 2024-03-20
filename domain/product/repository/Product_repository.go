package repository

import (
	"github.com/jmoiron/sqlx"
	"superindo/v1/domain/product"
)

type Product_repository struct {
	sqlx *sqlx.DB
}

func NewProduct_repository(sqlx *sqlx.DB) product.Repository {
	return &Product_repository{sqlx}
}

func (in *Product_repository) Rep_begin_transaction() (*sqlx.Tx, error) {
	tx, err := in.sqlx.Beginx()
	return tx, err
}

func (in *Product_repository) Rep_commit_transaction(tx *sqlx.Tx) error {

	err := tx.Commit()
	if err != nil {
		return in.Rep_rollback_transaction(tx)
	}
	return nil
}

func (in *Product_repository) Rep_rollback_transaction(tx *sqlx.Tx) error {
	var err error
	if err = tx.Rollback(); err != nil {
		return err
	}
	return err
}
