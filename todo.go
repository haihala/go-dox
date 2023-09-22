package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Key  int32
	Text string
}

type Todos struct {
	Value   []Todo
	Counter int32
	Mu      sync.Mutex
}

func (todos *Todos) Add(text string) int32 {
	todos.Mu.Lock()
	todos.Counter++
	todos.Value = append(todos.Value, Todo{Text: text, Key: todos.Counter})
	defer todos.Mu.Unlock()
	return todos.Counter
}

func (todos *Todos) Remove(id int32) {
	todos.Mu.Lock()
	for i, todo := range todos.Value {
		if todo.Key == id {
			todos.Value = append(todos.Value[:i], todos.Value[i+1:]...)
		}
	}
	todos.Mu.Unlock()
}

func (data *Todos) GetValue() map[string][]Todo {
	data.Mu.Lock()
	defer data.Mu.Unlock()
	return map[string][]Todo{
		"Todos": data.Value,
	}
}

func main() {
	todos := &Todos{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.ExecuteTemplate(w, "index.html", todos.GetValue())
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		field := r.Form.Get("new-todo")
		key := todos.Add(field)

		fmt.Fprintf(w, "<li id=\"todo-%d\">%s</li>", key, field)
	})

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
