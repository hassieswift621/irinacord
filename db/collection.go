package db

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/xerrors"
)

// Collection is the type for collection name consts.
type Collection string

// DeleteOne deletes a single document which matches the query.
func (c *Client) DeleteOne(ctx context.Context, col Collection, q interface{}) error {
	_, err := c.database.Collection(string(col)).DeleteOne(ctx, q)
	if err != nil {
		return xerrors.Errorf("delete document: %w", err)
	}

	return nil
}

// DeleteMany deletes multiple documents which match the query.
func (c *Client) DeleteMany(ctx context.Context, col Collection, q interface{}) error {
	_, err := c.database.Collection(string(col)).DeleteMany(ctx, q)
	if err != nil {
		return xerrors.Errorf("delete documents: %w", err)
	}

	return nil
}

// FindOne queries for a single document which matches the query.
// If a document is found, it is decoded into v.
func (c *Client) FindOne(ctx context.Context, col Collection, q interface{}, v interface{}) error {
	err := c.database.Collection(string(col)).FindOne(ctx, q).Decode(v)
	if err != nil {
		return xerrors.Errorf("find document: %w", err)
	}

	return nil
}

// FindMany queries for all documents which match the query.
// If documents are found, a cursor is returned.
func (c *Client) FindMany(ctx context.Context, col Collection, q interface{}) (*mongo.Cursor, error) {
	// Find documents.
	cursor, err := c.database.Collection(string(col)).Find(ctx, q)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, xerrors.Errorf("find documents: %w", err)
	}

	return cursor, nil
}

// FindManyReflect queries for all documents which match the query.
// If documents are found, they are decoded into v using reflection.
func (c *Client) FindManyReflect(ctx context.Context, col Collection, q interface{}, v interface{}) error {
	// Find documents.
	cursor, err := c.database.Collection(string(col)).Find(ctx, q)
	if err != nil {
		return xerrors.Errorf("find documents: %w", err)
	}

	// Decode documents through reflection.

	// v should be a slice.
	slice := reflect.ValueOf(v)
	if slice.Kind() != reflect.Ptr {
		return xerrors.New("v was not a pointer to a slice")
	}
	slice = slice.Elem()
	if slice.Kind() != reflect.Slice {
		return xerrors.New("v was not a slice")
	}

	// Get interface of slice element type.
	// If slice element is a pointer, i.e. []*E,
	// Elem() twice to get type, else Elem() once.
	var documentType reflect.Type
	if elemType := slice.Type(); elemType.Elem().Kind() == reflect.Ptr {
		documentType = elemType.Elem().Elem()
	} else {
		documentType = elemType.Elem()
	}

	// Decode documents.
	for cursor.Next(ctx) {
		document := reflect.New(documentType).Interface()
		err := cursor.Decode(document)
		if err != nil {
			return xerrors.Errorf("decode document: %w", err)
		}

		// Add document to slice through reflection.
		slice.Set(reflect.Append(slice, reflect.ValueOf(document)))
	}

	// Close cursor.
	err = cursor.Close(ctx)
	if err != nil {
		return xerrors.Errorf("close cursor: %w", err)
	}

	return nil
}

// InsertOne inserts a single document.
func (c *Client) InsertOne(ctx context.Context, col Collection, v interface{}) error {
	_, err := c.database.Collection(string(col)).InsertOne(ctx, v)
	if err != nil {
		return xerrors.Errorf("insert document: %w", err)
	}

	return nil
}

// UpsertOne upserts a single document.
func (c *Client) UpsertOne(ctx context.Context, col Collection, q interface{}, u interface{}) error {
	upsert := true
	_, err := c.database.Collection(string(col)).UpdateOne(ctx, q, u,
		&options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		return xerrors.Errorf("upsert document: %w", err)
	}

	return nil
}
