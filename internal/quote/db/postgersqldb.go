package db

import (
	"context"
	"fmt"

	"github.com/Kur04ka/telegram_bot/internal/quote"
	"github.com/jackc/pgx/v5"
)

var _ quote.Storage = &db{}

type db struct {
	connection *pgx.Conn
}

func (d *db) FindOne(ctx context.Context, id string) (quote.Quote, error) {
	var quote quote.Quote
	err := d.connection.QueryRow(context.Background(), "select * from quotes where id=$1", id).Scan(&quote.Id, &quote.Quote)
	if err != nil {
		return quote, fmt.Errorf("query row failed, error: %v", err)
	}

	return quote, nil
}

func NewStorage(connection *pgx.Conn) quote.Storage {
	return &db{
		connection: connection,
	}
}
