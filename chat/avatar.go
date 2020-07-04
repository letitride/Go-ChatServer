package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
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
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := filepath.Match(useridStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
