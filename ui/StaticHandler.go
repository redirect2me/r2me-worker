package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

//go:embed static
var embeddedFiles embed.FS

func maxAgeHandler(h http.Handler) http.Handler {
	maxAge, err := strconv.Atoi(os.Getenv("MAX_AGE"))
	if err != nil {
		maxAge = 3600 // 1 hour
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		h.ServeHTTP(w, r)
	})
}

func GetStaticHandler(logger *slog.Logger) http.Handler {

	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		logger.Error("unable to create static file system", "error", err)
		panic(err)
	}

	return maxAgeHandler(http.FileServer(http.FS(fsys)))
}
