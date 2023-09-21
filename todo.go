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

type Counter struct {
	value int
	mu    sync.Mutex
}

func (c *Counter) Increase() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *Counter) Decrease() {
	c.mu.Lock()
	c.value--
	c.mu.Unlock()
}

func (c *Counter) GetValue() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return map[string]int{
		"CounterValue": c.value,
	}
}

func RenderCounter(w http.ResponseWriter, data map[string]int) {
	tmplStr := "<div id=\"counter\">{{.CounterValue}}</div>"
	tmpl := template.Must(template.New("counter").Parse(tmplStr))
	tmpl.ExecuteTemplate(w, "counter", data)
}

func main() {
	counter := &Counter{}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("index.html")
		tmpl.ExecuteTemplate(w, "index.html", counter.GetValue())
	})

	r.Post("/increase", func(w http.ResponseWriter, r *http.Request) {
		counter.Increase()
		RenderCounter(w, counter.GetValue())
	})

	r.Post("/decrease", func(w http.ResponseWriter, r *http.Request) {
		counter.Decrease()
		RenderCounter(w, counter.GetValue())
	})

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8080", r))
}
