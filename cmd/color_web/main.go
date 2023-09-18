package main

import (
	_ "embed"
	"net/http"
	"os"
	"regexp"

	"github.com/actinium/openai-go-experiments/setup"
	"github.com/actinium/openai-go-experiments/tools/colorpicker"
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

		ctx := req.Context()
		prompt := req.URL.Query().Get("prompt")

		color, err := colorPicker.GenerateColor(ctx, prompt)
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
	setup.LoadEnv()

	chatClient, _, _ := setup.Clients()
	colorPicker := colorpicker.New(chatClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Get("/color", colorHandler(colorPicker))

	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}
