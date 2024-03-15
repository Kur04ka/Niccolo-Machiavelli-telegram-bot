package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func NewPostgresqlClient(user, password, host, port, dbname string) *pgx.Conn {
	postgresURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	connection, err := pgx.Connect(context.Background(), postgresURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database, error: %v\n", err)
		os.Exit(1)
	}

	return connection
}
