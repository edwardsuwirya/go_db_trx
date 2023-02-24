package main

import (
	"github.com/jmoiron/sqlx"
)

type TrxFn func(tx *sqlx.Tx) error

func WithStmtTransaction(db *sqlx.DB, fn TrxFn) error {
	tx := db.MustBegin()

	err := fn(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return rollbackErr
		}

		return err
	}

	return tx.Commit()
}
