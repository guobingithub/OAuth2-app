package main

import (
	"context"
	//"encoding/json"
	"OAuth2-demo/logger"
	"net/http"
	"OAuth2-demo/constants"

	"golang.org/x/oauth2"
	"fmt"
	//"io/ioutil"
	"io/ioutil"
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

	logger.Info("Client is running at 9001 port.")
	logger.Fatal(http.ListenAndServe(":9001", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	logger.Info(fmt.Sprintf("client loginHandler enter, Request:%v",r))
	u := config.AuthCodeURL(constants.AuthState)
	http.Redirect(w, r, u, http.StatusFound)
}

func callBackHandler(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hello callback!")
	logger.Info(fmt.Sprintf("callBackHandler enter, Request:%v",r))
	r.ParseForm()
	state := r.Form.Get("state")
	if state != constants.AuthState {
		logger.Error(fmt.Sprintf("callBackHandler, State invalid! state:%s",state))
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")
	if code == "" {
		logger.Error("callBackHandler, Code is null!")
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		logger.Error(fmt.Sprintf("callBackHandler, Exchange fail! err:%v, token:%v\n",err,token))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.Get("http://localhost:9000/auser?access_token=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	fmt.Fprintf(w, "111111111111111:%s\n", contents)

	//ex := make(map[string]interface{})
	//ex["words"] = "Hello callback!"
	//token.WithExtra(ex)
	//e := json.NewEncoder(w)
	//e.SetIndent("", "  ")
	//e.Encode(*token)
	//
	//logger.Info(fmt.Sprintf("callBackHandler ok, tokenInfo:%v\n",token))
}