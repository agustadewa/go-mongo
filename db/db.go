package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type post struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type db struct {
	ctx    context.Context
	client mongo.Client
	cancel context.CancelFunc
}

func (db *db) connect(uri string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	db.client = *client

	db.ctx, db.cancel = context.WithCancel(context.Background())
	err = db.client.Connect(db.ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func (db *db) deferThis() {
	db.cancel()
	db.client.Disconnect(db.ctx)
}

func (db *db) insertPost(title string, body string) {
	post := post{title, body}

	collection := db.client.Database("dev").Collection("posts")
	insertResult, err := collection.InsertOne(db.ctx, post)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func (db *db) getPost(title string) bson.M {
	post := bson.M{}

	collection := db.client.Database("dev").Collection("posts")
	err := collection.FindOne(db.ctx, bson.M{"title": title}).Decode(&post)

	if err != nil {
		fmt.Println(err)
	}

	return post
}

func (db *db) getPosts(title string) []bson.M {
	collection := db.client.Database("dev").Collection("posts")
	cursor, err := collection.Find(db.ctx, bson.M{"title": title})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(db.ctx)

	result := []bson.M{}

	for cursor.Next(db.ctx) {
		receive := bson.M{}
		err := cursor.Decode(&receive)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, receive)
	}
	return result
}

func (db *db) deletePost(title string) {
	collection := db.client.Database("dev").Collection("posts")
	delResult, err := collection.DeleteOne(db.ctx, bson.M{"title": title})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delResult)
}

func (db *db) deleteAll(title string) {
	collection := db.client.Database("dev").Collection("posts")
	delResult, err := collection.DeleteMany(db.ctx, bson.M{"title": title})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(delResult)
}
