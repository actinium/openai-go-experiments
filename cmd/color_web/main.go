package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/actinium/openai-go-experiments/setup"
	"github.com/actinium/openai-go-experiments/tools/colorpicker"
	"github.com/fatih/color"
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

		prompt := req.FormValue("prompt")
		if prompt == "" {
			http.Error(w, "Missing or empty prompt", 400)
			return
		}

		color, err := colorPicker.GenerateColor(req.Context(), prompt)
		if err != nil {
			http.Error(w, "Could not generate color", 500)
			return
		}

		if len(color) < 7 || !r.MatchString(color[0:7]) {
			http.Error(w, color, 500)
			return
		}

		w.Write([]byte(color[0:7]))
	}
}

func main() {
	setup.LoadEnv()

	chatClient := setup.Clients().Chat()
	colorPicker := colorpicker.New(chatClient)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Post("/color", colorHandler(colorPicker))

	link := color.New(color.Underline, color.FgHiBlue).SprintFunc()
	log.Printf("Listening on %s\n", link(os.Getenv("HTTP_ADDR")))
	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}
