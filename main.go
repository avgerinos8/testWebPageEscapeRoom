package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type PageData struct {
	Header          string
	SubHeader       string
	Text0           string
	Text1           string
	Text2           string
	CountNumber     int
	ContactMessage  string
	FormPlaceholder string
}

var data = PageData{
	Header:          "The Cryptic Grid: Escape Rooms",
	SubHeader:       "Escape from Reality",
	Text0:           "Έχεις μόνο 60 λεπτά για να λύσεις τους γρίφους. Ο χρόνος κυλάει",
	Text1:           "αντίστροφα.",
	Text2:           "Μπορείς να τα καταφέρεις;",
	CountNumber:     0,
	ContactMessage:  "Come and let us screw with your mind:",
	FormPlaceholder: "Let us have it...",
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

// ContactPage handles the form submission
func ContactPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read the text from the textarea input
	message := r.FormValue("comments")

	// Validate against empty messages
	if message == "" {
		data.ContactMessage = "You cannot send an empty soul into the grid."
		http.Redirect(w, r, "/#section3", http.StatusSeeOther)
		return
	}

	// Open or create contact.txt in append mode with write-only access
	file, err := os.OpenFile("contact.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close() // Ensure the file is closed at the end of the execution

	// Append the message followed by a new line character
	if _, err := file.WriteString(message + "\n"); err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Log the successful operation to the console
	fmt.Printf("Message saved successfully to contact.txt: %s\n", message)

	// Open HomePage after successful submission in section 3 with # with a success message
	data.ContactMessage = "Your thoughts have been trapped inside our mailbox..."
	data.FormPlaceholder = "Thank you for your contribution to the grid."
	http.Redirect(w, r, "/#section3", http.StatusSeeOther)
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/contact", ContactPage)
	mux.HandleFunc("/", HomePage)
	fmt.Println("Server ready. Use http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
