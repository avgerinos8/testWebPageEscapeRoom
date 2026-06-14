package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type PageData struct {
	Header      string
	SubHeader   string
	Text0       string
	Text1       string
	Text2       string
	CountNumber int
}

var data = PageData{
	Header:      "The Cryptic Grid: Escape Rooms",
	SubHeader:   "Escape from Reality",
	Text0:       "Έχεις μόνο 60 λεπτά για να λύσεις τους γρίφους. Ο χρόνος κυλάει",
	Text1:       "αντίστροφα.",
	Text2:       "Μπορείς να τα καταφέρεις;",
	CountNumber: 0,
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("http request GET with endpoint %s\n", r.URL.Path)
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	templ, err := template.ParseFiles("static/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data.CountNumber++
	templ.Execute(w, data)
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", HomePage)
	fmt.Println("Server ready. Use http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
