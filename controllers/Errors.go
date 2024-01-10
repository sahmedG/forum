package controllers

import (
	"html/template"
	"net/http"
)

type ErrorData struct {
	Err      int
	ImageURL string
}

func HTTPErrorHandler(w http.ResponseWriter, r *http.Request, errorCode int) {
	// Get the file name and line number of the function call
	// _, file, line, _ := runtime.Caller(1)
	// fmt.Printf("Error in %s:%d\n", file, line)
	w.WriteHeader(errorCode)
	errdata := ErrorData{Err: errorCode}
	// Set the status code based on the error code
	switch errorCode {
	case 500:
		errdata.ImageURL = "https://prod.spline.design/2SBPESoOzQ55Ovz7/scene.splinecode"
	case 400:
		errdata.ImageURL = "https://prod.spline.design/SnyQaOZHmrKTMfWz/scene.splinecode"
	case 404:
		errdata.ImageURL = "https://prod.spline.design/gG8xeIwuTZ1baGok/scene.splinecode"
	case 405:
		errdata.ImageURL = "https://media.tenor.com/HzKjCOw8gekAAAAd/baby-angry.gif"
	default:
		// If the error code is not recognized, return a generic error
		errdata.ImageURL = "https://prod.spline.design/2SBPESoOzQ55Ovz7/scene.splinecode"
	}

	// Execute the error template with the error data
	tmpl, err := template.ParseFiles("static/html/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, errdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
