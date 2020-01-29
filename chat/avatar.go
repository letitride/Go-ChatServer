package main

import (
	"errors"
)

//ErrNoAvatarURL はAvatarインスタンスがアバターURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLが取得できません。")

//Avatar はユーザのプロフィール画像を表す型です
type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

//AuthAvatar はauth認証情報からavatarを取得
type AuthAvatar struct{}

//UseAuthAvatar はJavaでいうstatic?　のように呼び出せる
var UseAuthAvatar AuthAvatar

//GetAvatarURL はclientのもつavatar_urlを返す
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}