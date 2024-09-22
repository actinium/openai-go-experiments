package npm

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:generate npm install --silent

//go:embed node_modules/alpinejs/dist/cdn.min.js
var alpinejs []byte

func AssetsHandler() http.Handler {
	r := chi.NewRouter()

	r.Get("/alpine.js", func(w http.ResponseWriter, req *http.Request) { w.Write(alpinejs) })

	return r
}
