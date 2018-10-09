package main

import (
	"OAuth2-demo/logger"
	"net/http"
	"net/url"
	"os"

	"github.com/go-session/session"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"OAuth2-demo/constants"
	."OAuth2-demo/oauth-flag"
	"fmt"
	"time"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set(constants.ClientId, &models.Client{
		ID:     constants.ClientId,
		Secret: constants.ClientSecret,
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewServer(server.NewConfig(), manager)

	srv.SetUserAuthorizationHandler(userAuthorizeHandler)
	srv.SetAccessTokenExpHandler(accessTokenExpHandler)
	srv.SetAuthorizeScopeHandler(authorizeScopeHandler)
	srv.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logger.Error("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		logger.Error("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/auser",func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("auser Handler enter, Request:%v",r))
		//exp,err := srv.AccessTokenExpHandler(w, r)
		//if err != nil {
		//	logger.Error(fmt.Sprintf("auser Handler, AccessTokenExpHandler err:%v",err))
		//	http.Error(w, err.Error(), http.StatusBadRequest)
		//}
		//
		//logger.Error(fmt.Sprintf("auser Handler, exp:%v",exp))

		r.ParseForm()
		accessToken := r.Form.Get("access_token")
		logger.Error("auser Handler, request form access_token:",accessToken)
		ti,err := srv.Manager.LoadAccessToken(accessToken)
		if err!=nil{
			logger.Error(fmt.Sprintf("auser Handler, accessToken fail! err:%v",err))
			return
		}
		logger.Error(fmt.Sprintf("auser Handler, accessToken ok! ti:%v",ti))
		logger.Error("auser Handler, accessToken ok! userId:",ti.GetUserID())
		logger.Error("auser Handler, accessToken ok! access_token:",ti.GetAccess())

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello guoBin@."))
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("authorize Handler enter, Request:%v",r))
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			logger.Error(fmt.Sprintf("authorize Handler, HandleAuthorizeRequest err:%v",err))
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("token Handler enter, Request:%v",r))
		err := srv.HandleTokenRequest(w, r)
		if err != nil {
			logger.Error(fmt.Sprintf("token Handler, HandleTokenRequest err:%v",err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		logger.Error(fmt.Sprintf("token Handler, HandleTokenRequest ok, w:%v\n, r:%v",w,r))
	})

	logger.Info("Server is running at 9000 port.")
	logger.Fatal(http.ListenAndServe(":9000", nil))
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	logger.Info(fmt.Sprintf("userAuthorizeHandler enter, Request:%v",r))
	store, err := session.Start(nil, w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("UserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		logger.Info("userAuthorizeHandler, ReturnUri:",r.Form)
		logger.Info("userAuthorizeHandler, has no userId, will login.")
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("UserID")
	store.Save()
	return
}

func accessTokenExpHandler(w http.ResponseWriter, r *http.Request) (exp time.Duration, err error){
	logger.Info(fmt.Sprintf("accessTokenExpHandler enter, Request:%v",r))
	return 120 * time.Second, nil
}

func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error){
	logger.Info(fmt.Sprintf("authorizeScopeHandler enter, Request:%v",r))
	return "user", nil
}

func passwordAuthorizationHandler(username, password string) (userID string, err error){
	if username == "admin" && password == "123456" {
		userID = "000000"
		return
	}
	err = fmt.Errorf("user not found")
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("loginHandler enter, Request:%v",r))
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Error("1111111111111111111111111")
	sid := store.SessionID()
	logger.Error("sssssssssssssssssssssid:",sid)
	if r.Method == "POST" {
		if r.Form == nil {
			r.ParseForm()
		}
		logger.Error("login user:",r.Form.Get("user"))

		logger.Info("loginHandler post")
		store.Set("LoggedInUserID", r.Form.Get("user"))
		store.Save()

		w.Header().Set("Location", "/auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	showHTML(w, r, LoginPath)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("authHandler enter, Request:%v",r))
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		logger.Error("authHandler, LoggedInUserID not found, will re login.")
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	if r.Method == "POST" {
		logger.Error("authHandler, post.")
		var form url.Values
		if v, ok := store.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		u := new(url.URL)
		u.Path = "/authorize"
		u.RawQuery = form.Encode()
		w.Header().Set("Location", u.String())
		logger.Error("uuuuuuuuuuuu Location:",u.String())
		w.WriteHeader(http.StatusFound)
		store.Delete("Form")

		if v, ok := store.Get("LoggedInUserID"); ok {
			store.Set("UserID", v)
		}
		store.Save()

		return
	}

	showHTML(w, r, AuthPath)
}

func showHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

func userInfoHandler(w http.ResponseWriter, r *http.Request){
	logger.Info(fmt.Sprintf("userInfoHandler enter, Request:%v",r))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello guoBin."))
}
