package main

//gomuniauth/commonのaliasとしてgomniauthcommonとする
import gomniauthcommon "github.com/stretchr/gomniauth/common"
import (
	"crypto/md5"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"io"
	"log"
	"net/http"
	"strings"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	//type embedding
	gomniauthcommon.User //AvatarURL()はUserインターフェースで宣言済み
	uniqueID             string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

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
		//認証プロバイダページへリダイレクト
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}

		//r.URL.RawQuery GETパラメータ全体文字列
		//callback時のproviderから渡されたパラメータから認証を完了させる
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした:", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました:", provider, "-", err)
		}

		chatUser := &chatUser{User: user}
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("GetAvatarURLに失敗しました", "-", err)
		}
		//user名をname:usernameとしてbase64encode
		authCookieValue := objx.New(map[string]interface{}{
			"userid":     chatUser.uniqueID,
			"name":       user.Name(),
			"avatar_url": avatarURL,
		}).MustBase64()
		//base64した文字列をauthというキー名でcookieに保存
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション %s には非対応です", action)
	}
}
