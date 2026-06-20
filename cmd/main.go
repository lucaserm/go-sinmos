package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/lucaserm/go-sinmos/internal/env"
)

func main() {
	cfg := config{
		addr: ":8080",
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", "host=localhost user=postgres password=postgres dbname=sinmos sslmode=disable"),
		},
	}

	//logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	logger.Info("connected to database")

	// api
	app := application{
		config: cfg,
		db:     conn,
	}

	router := app.mount()
	router = app.v1(router)
	if err := app.run(router); err != nil {
		slog.Error("server has failed to start", "error", err)
		os.Exit(1)
	}
}
