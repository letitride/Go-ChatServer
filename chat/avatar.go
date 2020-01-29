package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
func (a AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

//GravatarAvatar はgravatorからavator情報を取得
type GravatarAvatar struct{}

//UseGravatar はクラスメソッドアクセサ
var UseGravatar GravatarAvatar

//GetAvatarURL はGravatorAvatarのinterface実装
func (a GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
