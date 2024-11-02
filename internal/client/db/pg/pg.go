package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/Mobo140/microservices/chat/internal/client/db"
	"github.com/Mobo140/microservices/chat/internal/client/db/prettier"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ db.DB = (*pg)(nil)

type pg struct {
	dbc *pgxpool.Pool
}

func NewDB(dbc *pgxpool.Pool) *pg {
	return &pg{dbc: dbc}
}

func (p *pg) Close() {
	p.dbc.Close()
}

func (p *pg) Ping(ctx context.Context) error {
	conn, err := p.dbc.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return conn.Conn().Ping(ctx)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...interface{}) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	return p.dbc.Exec(ctx, q.QueryRow, args...)
}

func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...interface{}) (pgx.Rows, error) {

	return p.dbc.Query(ctx, q.QueryRow, args...)
}

func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...interface{}) pgx.Row {
	logQuery(ctx, q, args...)

	return p.dbc.QueryRow(ctx, q.QueryRow, args...)
}

func (p *pg) ScanOneContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest interface{}, q db.Query, args ...interface{}) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func logQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRow, prettier.PlaceHolderDollar, args...)
	log.Println(
		ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
