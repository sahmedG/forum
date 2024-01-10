package api

import (
	"encoding/json"
	"fmt"
	"forum/pkgs/funcs"
	"net/http"
)

func Serve_categories_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	container, err := funcs.GetAllCategoryInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get category information: %v", err), http.StatusInternalServerError)
		return
	}

	/* Marshal the data into JSON format */
	jsonData, err := json.Marshal(container)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling JSON: %v", err), http.StatusInternalServerError)
		return
	}

	/* Set the content type to JSON */
	w.Header().Set("Content-Type", "application/json")

	/* Write the JSON data to the response writer */
	w.Write(jsonData)
}
