package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/authorize",
			TokenURL: "http://localhost:8080/token",
		},
	}
)

// Homepage
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Homepage Hit!")
	u := config.AuthCodeURL("xyz")
	http.Redirect(w, r, u, http.StatusFound)
}

// Authorize
func Authorize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
}

func main() {
	http.HandleFunc("/", HomePage)

	http.HandleFunc("/oauth2", Authorize)

	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
