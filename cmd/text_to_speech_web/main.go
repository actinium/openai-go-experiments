package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/actinium/openai-go-experiments/npm"
	"github.com/actinium/openai-go-experiments/openai"
	"github.com/actinium/openai-go-experiments/setup"
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

func textToSpeechHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		text := req.FormValue("text")
		if text == "" {
			http.Error(w, "Missing or empty text", 400)
			return
		}

		voice := req.FormValue("voice")
		if voice == "" {
			http.Error(w, "Missing or empty voice", 400)
			return
		}
		if !slices.Contains(openai.TextToSpeechVoices[:], voice) {
			http.Error(w, "Voice is not valid", 400)
			return
		}

		tts := setup.Clients().TTS()
		tts.SetHD(true)
		tts.SetVoice(voice)

		audio, err := tts.GenerateAudioWithContext(req.Context(), text)
		if err != nil {
			http.Error(w, "Could not generate audio", 500)
			return
		}

		w.Header().Add("content-type", "audio/mpeg")
		w.Write(audio.Data)
	}
}

func main() {
	setup.LoadEnv()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/assets/*", assetsHandler())

	r.Get("/", indexHandler())
	r.Post("/tts", textToSpeechHandler())

	r.Mount("/npm", npm.AssetsHandler())

	link := color.New(color.Underline, color.FgHiBlue).SprintFunc()
	log.Printf("Listening on %s\n", link(os.Getenv("HTTP_ADDR")))
	http.ListenAndServe(os.Getenv("HTTP_ADDR"), r)
}
