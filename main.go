package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/beldmian/go-markdown-server/db"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

// Post ...
type Post struct {
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
	collectionResp, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	collection = collectionResp
	r := mux.NewRouter()
	configureRouter(r)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func configureRouter(r *mux.Router) {
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/post/{name}", mdNamedHandler)
	r.HandleFunc("/add", addHandler)
}

func errorNotFoundPage(w http.ResponseWriter) {
	text := `# Error 404
This page not found`
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(text))))
	tmpl.ExecuteTemplate(w, "md", output)
}

func internalServerErrorPage(err error, w http.ResponseWriter) {
	log.Panic(err)
	text := `# Error 500
Internal server error`
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(text))))
	tmpl.ExecuteTemplate(w, "md", output)
}
