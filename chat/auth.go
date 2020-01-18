package main

import "net/http"

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
