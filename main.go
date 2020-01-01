package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/russross/blackfriday.v2"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/md/{name}", mdNamedHandler)
	fmt.Println("Server have started")
	http.ListenAndServe(":3000", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("md/index.md")
	tmpl := template.Must(template.ParseFiles("md.html"))
	output := template.HTML(string(blackfriday.Run(file)))
	tmpl.ExecuteTemplate(w, "md", output)
}
func mdNamedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, err := ioutil.ReadFile("md/" + vars["name"] + ".md")
	if err != nil {
		fmt.Fprint(w, "File Not found")
	} else {
		tmpl := template.Must(template.ParseFiles("md.html"))
		output := template.HTML(string(blackfriday.Run(file)))
		tmpl.ExecuteTemplate(w, "md", output)
	}
}
