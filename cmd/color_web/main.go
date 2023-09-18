package main

import (
	_ "embed"
	"net/http"
	"openaigo/setup"
	"openaigo/tools/colorpicker"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed index.html
var page []byte

func indexHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(page)
	}
}

func colorHandler(colorPicker *colorpicker.ColorPicker) func(http.ResponseWriter, *http.Request) {

	r, _ := regexp.Compile("^#[0-9A-Fa-f]+$")

	return func(w http.ResponseWriter, req *http.Request) {
		if !req.URL.Query().Has("prompt") {
			http.Error(w, "Missing prompt", 400)
			return
		}
		color, err := colorPicker.Color(req.URL.Query().Get("prompt"))
		if err != nil {
			http.Error(w, "Could not generate color", 500)
			return
		}

		if !r.MatchString(color[0:7]) {
			http.Error(w, color, 500)
			return
		}

		w.Write([]byte(color[0:7]))
	}
}

func main() {
	chatClient, _ := setup.Clients()
	colorPicker := colorpicker.New(chatClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Get("/color", colorHandler(colorPicker))

	http.ListenAndServe("localhost:8090", r)
}
