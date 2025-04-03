package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Name    string
	Content string
}

var messages []Message

func loadMessagesFromFile() {
	file, err := os.Open("messages.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Println("error reading messages file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			messages = append(messages, Message{
				Name:    parts[0],
				Content: parts[1],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Print("Error scanning file", err)
	}
}
func main() {
	loadMessagesFromFile() // loads saved input messages first

	http.HandleFunc("/", handler)
	http.HandleFunc("/submit", submitHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Template parsing error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, messages)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.FormValue("name")
	content := r.FormValue("message")

	newMessage := Message{Name: name, Content: content}
	messages = append(messages, newMessage)

	// saves user input to .txt file
	f, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("File error:", err)
	} else {
		defer f.Close()
		_, err := f.WriteString(fmt.Sprintf("%s: %s\n", name, content))
		if err != nil {
			log.Println("Write error:", err)
		}
	}

	// redirect back to homepage after user submits
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
