package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

// CreateRequest ...
type CreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
}

var (
	port       string
	collection *mongo.Collection
)

func init() {
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	} else {
		port = ":8080"
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://beld:124252@cluster0-wmuco.mongodb.net/blog"))
	collection = client.Database("blog").Collection("posts")
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/post/{name}", mdNamedHandler)
	r.HandleFunc("/add", addHandler)
	fmt.Println("Server have started")
	log.Print("Server started", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	out := "# Home\n---\n"
	ctx5, cancel5 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel5()
	for cursor.Next(ctx5) {
		var post CreateRequest
		err := cursor.Decode(&post)
		if err != nil {
			log.Fatal(err)
		}
		out += "- [" + post.Title + "](/post/" + post.URL + ")\n"
	}
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(out))))
	if err := tmpl.ExecuteTemplate(w, "md", output); err != nil {
		log.Fatal(err)
	}
}
func mdNamedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"url": vars["name"]}
	var post CreateRequest
	if err := collection.FindOne(ctx, filter).Decode(&post); err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(post.Body))))
	tmpl.ExecuteTemplate(w, "md", output)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	//if r.FormValue("key") == "124252" {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := CreateRequest{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
		URL:   r.FormValue("url"),
	}
	obj, err := collection.InsertOne(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Add object with id", obj.InsertedID)
	w.Write([]byte("add success"))
	//}
}
