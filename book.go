package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Book struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Author string             `json:"author,omitempty" bson:"author,omitempty"`
}

var bookCollectionName = "books"

func getBooks(db *mongo.Database, start, count int) ([]Book, error) {
	col := booksCollection(db)
	ctx := dbContext(30)

	cursor, err := col.Find(ctx, bson.M{}) // find all
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bs []Book
	for cursor.Next(ctx) {
		var b Book
		cursor.Decode(&b)
		bs = append(bs, b)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bs, nil
}

func (b *Book) getBook(db *mongo.Database) error {
	col := booksCollection(db)
	ctx := dbContext(30)

	filter := bson.M{"_id": b.ID}
	err := col.FindOne(ctx, filter).Decode(&b)
	return err
}

func (b *Book) createBook(db *mongo.Database) (*mongo.InsertOneResult, error) {
	col := booksCollection(db)
	ctx := dbContext(30)

	result, err := col.InsertOne(ctx, b)
	// Convert to map[string]string
	// id := map[string]string{"_id": result.InsertedID.(primitive.ObjectID).Hex()}
	return result, err
}

func (b *Book) updateBook(db *mongo.Database, ub Book) (*mongo.UpdateResult, error) {
	col := booksCollection(db)
	ctx := dbContext(30)

	filter := bson.M{"_id": b.ID}
	update := bson.M{"$set": ub}
	result, err := col.UpdateOne(ctx, filter, update)
	// Get the updated document
	// col.FindOne(ctx, filter).Decode(&b)
	return result, err
}

func (b *Book) deleteBook(db *mongo.Database) (*mongo.DeleteResult, error) {
	col := booksCollection(db)
	ctx := dbContext(30)

	filter := bson.M{"_id": b.ID}
	result, err := col.DeleteOne(ctx, filter)
	return result, err
}

// helpers
func booksCollection(db *mongo.Database) *mongo.Collection {
	return db.Collection(bookCollectionName)
}

func dbContext(i time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), i*time.Second)
	return ctx
}
