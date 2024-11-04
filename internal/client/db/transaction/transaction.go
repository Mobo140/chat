package transaction

import (
	"context"

	"github.com/Mobo140/microservices/chat/internal/client/db"
	"github.com/Mobo140/microservices/chat/internal/client/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewManager(db db.Transactor) db.TxManager {
	return &manager{db: db}
}

func (m *manager) ReadCommited(ctx context.Context, fn db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, fn)
}

func (m *manager) transaction(ctx context.Context, txOpts pgx.TxOptions, fn db.Handler) error {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := m.db.BeginTx(ctx, txOpts)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	ctx = pg.MakeContext(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "failed to commit transaction")
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}
