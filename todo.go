package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Key  int
	Text string
}

type Todos struct {
	Items   []Todo
	Counter int
	Mu      sync.Mutex
}

func (todos *Todos) Add(text string) Todo {
	todos.Mu.Lock()
	todos.Counter++
	todos.Items = append(todos.Items, Todo{Text: text, Key: todos.Counter})
	defer todos.Mu.Unlock()
	return todos.Items[len(todos.Items)-1]
}

func (todos *Todos) Remove(id int) {
	todos.Mu.Lock()
	for i, todo := range todos.Items {
		if todo.Key == id {
			todos.Items = append(todos.Items[:i], todos.Items[i+1:]...)
		}
	}
	todos.Mu.Unlock()
}

func (data *Todos) GetValue() map[string][]Todo {
	data.Mu.Lock()
	defer data.Mu.Unlock()
	return map[string][]Todo{
		"Todos": data.Items,
	}
}

func main() {
	todos := &Todos{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("index.html", "todo.html")
		tmpl.ExecuteTemplate(w, "index.html", todos.GetValue())
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		field := r.Form.Get("new-todo")
		todo := todos.Add(field)

		tmpl, _ := template.ParseFiles("todo.html")
		tmpl.ExecuteTemplate(w, "todo", todo)
	})

	r.Post("/remove/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		todos.Remove(id)
	})

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
