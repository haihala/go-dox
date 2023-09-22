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
	Value   []Todo
	Counter int
	Mu      sync.Mutex
}

func (todos *Todos) Add(text string) int {
	todos.Mu.Lock()
	todos.Counter++
	todos.Value = append(todos.Value, Todo{Text: text, Key: todos.Counter})
	defer todos.Mu.Unlock()
	return todos.Counter
}

func (todos *Todos) Remove(id int) {
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
		tmpl.Execute(w, todos.GetValue())
	})

	r.Post("/add", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		field := r.Form.Get("new-todo")
		key := todos.Add(field)
		template_data := map[string]string{
			"Id":      fmt.Sprint(key),
			"Content": field,
		}

		tmpl, _ := template.ParseFiles("test.html")
		tmpl.Execute(w, template_data)
	})

	r.Post("/remove/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		todos.Remove(id)
	})

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
