package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/xerrors"
)

const db = "dbIrinacord"

type Client struct {
	client   *mongo.Client
	database *mongo.Database
}

// New creates a new instance of the MongoDB client.
func New(config *ClientConfig) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, xerrors.Errorf("create mongo client: %w", err)
	}

	return &Client{client: client}, nil
}

// Connect connects the client to the MongoDB server.
func (c *Client) Connect(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return xerrors.Errorf("connect mongo client: %w", err)
	}

	return nil
}

// Disconnect disconnects the client from the MongoDB server.
func (c *Client) Disconnect(ctx context.Context) error {
	err := c.client.Disconnect(ctx)
	if err != nil {
		return xerrors.Errorf("disconnect mongo client: %w", err)
	}

	return nil
}

// Ping pings the MongoDB server.
// If no error occurred, a handle to the database is stored in the client.
func (c *Client) Ping(ctx context.Context) error {
	err := c.client.Ping(ctx, readpref.Primary())
	if err != nil {
		c.database = nil
		return xerrors.Errorf("ping mongo: %w", err)
	}

	c.database = c.client.Database(db)

	return nil
}
