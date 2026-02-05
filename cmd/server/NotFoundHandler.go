package main

import (
	"net/http"

	"github.com/redirect2me/r2me-worker/ui"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	content, err := ui.ExpandTemplate("templates/404.gohtml", map[string]interface{}{
		"RequestPath": r.URL.Path,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Error("unable to expand template", "error", err)
		return
	}

	w.Write([]byte(content))
}
