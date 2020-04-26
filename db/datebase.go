package db

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Post ...
type Post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

// ConnectToDB ...
func ConnectToDB() (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://beld:124252@cluster0-wmuco.mongodb.net/blog"))
	if err != nil {
		return nil, err
	}
	collection := client.Database("blog").Collection("posts")
	return collection, nil
}

// GetPosts ...
func GetPosts(collection *mongo.Collection) ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "_id", Value: -1}})
	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return []Post{}, err
	}
	ctx5, cancel5 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel5()
	var posts []Post
	for cursor.Next(ctx5) {
		var post Post
		err := cursor.Decode(&post)
		if err != nil {
			return []Post{}, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// GetPostByName ...
func GetPostByName(collection *mongo.Collection, name string) (Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"url": name}
	var post Post
	if err := collection.FindOne(ctx, filter).Decode(&post); err != nil {
		return Post{}, err
	}
	return post, nil
}

// InsertPost ...
func InsertPost(collection *mongo.Collection, post Post, key string) (*mongo.InsertOneResult, error) {
	if key == "124252" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		obj, err := collection.InsertOne(ctx, post)
		if err != nil {
			return nil, err
		}
		return obj, err
	}
	return nil, errors.New("Key is not valid")
}
