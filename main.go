package main

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type User struct {
	Username string
	Password string
}

var users map[string]string

func loadUsers(filePath string) {
	users = make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		if len(record) == 2 {
			users[record[0]] = record[1]
		}
	}
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	password := hashPassword(r.FormValue("password"))

	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Hashed Password: %s\n", password)
	fmt.Printf("Expected Hash: %s\n", users[username])

	if storedPassword, exists := users[username]; exists && storedPassword == password {
		switch username {
		case "hackerman":
			http.Redirect(w, r, "/hackerman", http.StatusSeeOther)
		case "doingus":
			http.Redirect(w, r, "/architecture", http.StatusSeeOther)
		default:
			http.Redirect(w, r, "/weird", http.StatusSeeOther)
		}
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func hackermanPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/hackerman.html")
	tmpl.Execute(w, nil)
}

func weirdPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/weird.html")
	tmpl.Execute(w, nil)
}

func architecturePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/architecture.html")
	tmpl.Execute(w, nil)
}

func main() {
	loadUsers("users.secure")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/hackerman", hackermanPage)
	http.HandleFunc("/weird", weirdPage)
	http.HandleFunc("/architecture", architecturePage)

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
