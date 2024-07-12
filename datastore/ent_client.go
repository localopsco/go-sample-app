package datastore

import (
	"context"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/localopsco/go-sample/ent"
)

func NewEntClient(host, port, userName, password, dbName string) (*ent.Client, error) {
	client, err := ent.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host,
			port,
			userName,
			dbName,
			password,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}
	return client, nil
}
