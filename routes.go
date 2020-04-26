package main

import (
	"html/template"
	"net/http"

	"github.com/beldmian/go-markdown-server/db"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := db.GetPosts(collection)
	if err != nil {
		internalServerErrorPage(err, w)
	}
	out := "# Home\n---\n"
	for _, post := range posts {
		out += "- [" + post.Title + "](/post/" + post.URL + ")\n"
	}
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(out))))
	if err := tmpl.ExecuteTemplate(w, "md", output); err != nil {
		internalServerErrorPage(err, w)
	}
}
func mdNamedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, err := db.GetPostByName(collection, vars["name"])
	if err != nil {
		errorNotFoundPage(w)
	}
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run([]byte(post.Body))))
	tmpl.ExecuteTemplate(w, "md", output)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	post := db.Post{
		Title: v["title"],
		URL:   v["url"],
		Body:  v["body"],
	}
	_, err := db.InsertPost(collection, post, v["key"])
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("success"))
}
