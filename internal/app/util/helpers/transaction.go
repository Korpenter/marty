package helpers

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CommitTx(ctx context.Context, tx pgx.Tx, err error) {
	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}
}
