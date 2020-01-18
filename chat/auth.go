package main

import "net/http"

import "strings"

import "log"

import "fmt"

import "github.com/stretchr/gomniauth"

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//認証クッキーの取得
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		//何らかのエラー
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

//MustAuth 認証確認の為の事前処理handleを返します
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

//パスの形式 /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", provider, "-", err)
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション %s には非対応です", action)
	}
}
