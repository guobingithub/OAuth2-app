package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"OAuth2-demo/constants"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     constants.ClientId,
		ClientSecret: constants.ClientSecret,
		Scopes:       []string{"user"},
		RedirectURL:  constants.DomainUrl + "/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  constants.AuthURL,
			TokenURL: constants.TokenURL,
		},
	}
)

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callBackHandler)

	log.Println("Client is running at 9001 port.")
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	u := config.AuthCodeURL(constants.AuthState)
	http.Redirect(w, r, u, http.StatusFound)
}

func callBackHandler(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hello callback!")
	r.ParseForm()
	state := r.Form.Get("state")
	if state != constants.AuthState {
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
	ex := make(map[string]interface{})
	ex["words"] = "Hello callback!"
	token.WithExtra(ex)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
}