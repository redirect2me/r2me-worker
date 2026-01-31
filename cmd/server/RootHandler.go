package main

import (
	"net/http"

	"github.com/redirect2me/r2me-worker/ui"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	content, err := ui.ExpandTemplate("templates/index.gohtml", nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		Logger.Error("unable to expand template", "error", err)
		return
	}

	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(content))
}
