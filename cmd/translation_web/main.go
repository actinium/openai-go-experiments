package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/actinium/openai-go-experiments/setup"
	"github.com/actinium/openai-go-experiments/tools/translator"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed index.html
var page []byte

//go:embed assets
var assets embed.FS

func indexHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(page)
	}
}

func assetsHandler() func(http.ResponseWriter, *http.Request) {
	isDir := func(path string) bool {
		return strings.HasSuffix(path, "/")
	}

	return func(w http.ResponseWriter, req *http.Request) {
		if isDir(req.URL.Path) {
			http.NotFound(w, req)
			return
		}

		http.FileServer(http.FS(assets)).ServeHTTP(w, req)
	}
}

func translationHandler(translator *translator.Translator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fromLanguage := req.FormValue("from_language")
		if fromLanguage == "" {
			http.Error(w, "Missing or empty parameter: from_language", 400)
			return
		}
		toLanguage := req.FormValue("to_language")
		if toLanguage == "" {
			http.Error(w, "Missing or empty parameter: to_language", 400)
			return
		}
		text := req.FormValue("text")
		if text == "" {
			http.Error(w, "Missing or empty parameter: text", 400)
			return
		}

		translation, err := translator.Translate(
			req.Context(),
			fromLanguage,
			toLanguage,
			text,
		)
		if err != nil {
			http.Error(w, "Could not generate color", 500)
			return
		}

		w.Write([]byte(translation))
	}
}

func main() {
	setup.LoadEnv()

	translator := translator.New(setup.ChatClient())

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", indexHandler())
	r.Get("/assets/*", assetsHandler())
	r.Post("/translate", translationHandler(translator))

	link := color.New(color.Underline, color.FgHiBlue).SprintFunc()
	log.Printf("Listening on %s\n", link(os.Getenv("HTTP_ADDR")))
	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}
