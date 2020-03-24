package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

// CreateRequest ...
type CreateRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/post/{name}", mdNamedHandler)
	r.HandleFunc("/add", addHandler)
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

func addHandler(w http.ResponseWriter, r *http.Request) {
	req := CreateRequest{
		Title: r.FormValue("Title"),
		Body:  r.FormValue("Body"),
	}
	f, err := os.Create("./md/" + req.Title + ".md")
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte(req.Body))
	defer f.Close()
	index, err := os.OpenFile("./md/index.md", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = index.WriteString("\n- [" + req.Title + "](/post/" + req.Title + ")"); err != nil {
		log.Fatal(err)
	}
	defer index.Close()
}
