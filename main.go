package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"unicode"
)

type Message struct {
	Value    string `json:"value"`
	Editable bool   `json:"editable"`
}

type Amount struct {
	Value    int  `json:"value"`
	Editable bool `json:"editable"`
}

type Config struct {
	Size    int     `json:"size"`
	Border  int     `json:"border"`
	Payee   string  `json:"payee"`
	Color   bool    `json:"color"`
	Message Message `json:"message"`
	Amount  Amount  `json:"amount"`
}

func main() {
	intFn := func(s string) (any, error) {
		return strconv.Atoi(s)
	}

	boolFn := func(s string) (any, error) {
		return strconv.ParseBool(s)
	}

	parameters := map[string]func(s string) (any, error){
		"size":             intFn,
		"border":           intFn,
		"color":            boolFn,
		"amount":           intFn,
		"amout_editable":   boolFn,
		"message":          func(s string) (any, error) { return s, nil },
		"message_editable": boolFn,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:] // Strip the leading slash

		if path == "" || !isNumeric(path) {
			http.Error(w, "Invalid path parameter", http.StatusBadRequest)
			return
		}

		request := Config{
			Size:    240,
			Border:  1,
			Payee:   path,
			Color:   true,
			Message: Message{Editable: true},
			Amount:  Amount{Editable: true},
		}

		for param, fn := range parameters {
			paramValue := r.URL.Query().Get(param)
			if paramValue == "" {
				continue
			}

			value, err := fn(paramValue)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid query parameter for '%s': %v", param, err), http.StatusBadRequest)
				return
			}

			switch param {
			case "size":
				request.Size = value.(int)
			case "border":
				request.Border = value.(int)
			case "color":
				request.Color = value.(bool)
			case "amount":
				request.Amount.Value = value.(int)
			case "amount_editable":
				request.Amount.Editable = value.(bool)
			case "message":
				request.Message.Value = value.(string)
			case "message_editable":
				request.Message.Editable = value.(bool)
			}
		}

		j, err := json.Marshal(request)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create request body: %v", err), http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequest("POST", "https://api.swish.nu/qr/v2/prefilled", bytes.NewBuffer(j))
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		// Forward the original request headers to the backend request
		r.Header.Add("Content-Type", "application/json")
		for key, values := range r.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch content from backend: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		// Copy the content from the backend response to the proxy response
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to forward response content: %v", err), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func isNumeric(s string) bool {
	for _, c := range s {
		if !unicode.IsNumber(c) {
			return false
		}
	}

	return true
}
