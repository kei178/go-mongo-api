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
	col := db.Collection(bookCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

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
	col := db.Collection(bookCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	filter := bson.M{"_id": b.ID}
	err := col.FindOne(ctx, filter).Decode(&b)
	return err
}

func (b *Book) createBook(db *mongo.Database) (map[string]string, error) {
	col := db.Collection(bookCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	result, err := col.InsertOne(ctx, b)
	id := map[string]string{"_id": result.InsertedID.(primitive.ObjectID).Hex()}
	return id, err
}

func (b *Book) updateBook(db *mongo.Database, ub Book) error {
	col := db.Collection(bookCollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	filter := bson.M{"_id": b.ID}
	update := bson.M{"$set": ub}
	_, err := col.UpdateOne(ctx, filter, update)
	// Get the updated document
	col.FindOne(ctx, filter).Decode(&b)
	return err
}
