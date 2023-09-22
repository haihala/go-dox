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
	Text string
}

type Todos struct {
	Value []Todo
	Mu    sync.Mutex
}

func (todos *Todos) Add(todo Todo) {
	todos.Mu.Lock()
	todos.Value = append(todos.Value, todo)
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
	todos := &Todos{
		Value: []Todo{
			{Text: "Test1"},
			{Text: "Test2"},
		},
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.ExecuteTemplate(w, "index.html", todos.GetValue())
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		field := r.Form.Get("new-todo")
		todos.Add(Todo{Text: field})

		fmt.Fprintf(w, "<li>%s</li>", field)
	})

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
