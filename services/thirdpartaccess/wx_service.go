package thirdpartaccess

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

type WxConfig struct {
	Appid     string // 小程序配置文件
	AppSecret string
	ApiAddr   string
}

type WxLogin struct {
	appid     string // 小程序配置文件
	appSecret string

	SessionKey string
	UnionId    string // 小程序配置文件
	OpenID     string

	AccessToken string
	ExpiresIn   time.Time //Token过期时间
	wxManager   *Manager  //http客户端
}

func NewWxLoginService(appid, appsecret, apiaddr string) (*WxLogin, error) {
	if len(appid) == 0 || len(appsecret) == 0 || len(apiaddr) == 0 {
		return nil, errors.New("empty param")
	}
	wxManager := NewManager(apiaddr)
	wxLogin := &WxLogin{
		appid:     appid,
		appSecret: appsecret,
		wxManager: wxManager,
	}
	return wxLogin, nil
}

/*
	查询openid unionid 等信息
*/
func (w *WxLogin) QuerySession(code string) (*WxSession, error) {

	var ret WxSession

	query := url.Values{}
	setQuery(query, "appid", w.appid)
	setQuery(query, "secret", w.appSecret)
	setQuery(query, "js_code", code)
	setQuery(query, "grant_type", "authorization_code")
	logs.Debug(query.Encode())
	err := w.wxManager.client.Call(context.Background(), &ret, "GET", w.wxManager.url("/sns/jscode2session?%v", query.Encode()), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	获取接口调用凭据
*/
func (w *WxLogin) FreshAccessToken() error {
	token, err := w.GetAccessToken()
	if err != nil {
		logs.Warn("WX Service Fresh Access Token Failed!")
		return err
	}
	w.AccessToken = token.AccessToken
	w.ExpiresIn = time.Now().Add(time.Duration(token.ExpiresIn))
	return nil
}

/*
	获取接口调用凭据
*/
func (w *WxLogin) GetAccessToken() (*WxToken, error) {

	var ret WxToken

	query := url.Values{}
	setQuery(query, "appid", w.appid)
	setQuery(query, "secret", w.appSecret)
	setQuery(query, "grant_type", "client_credential")

	err := w.wxManager.client.Call(context.Background(), &ret, "GET", w.wxManager.url("/cgi-bin/token?%v", query.Encode()), nil)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

/*
	获取接口调用凭据
*/
func (w *WxLogin) GetPhoneNumber(code string) (*WxPhoneNum, error) {

	// 判断accesstoken是否过期
	if time.Now().After(w.ExpiresIn) {
		if err := w.FreshAccessToken(); err != nil {
			return nil, err
		}
	}

	query := url.Values{}
	setQuery(query, "access_token", w.AccessToken)

	codeJson := struct {
		Code string `json:"code"`
	}{Code: code}
	var ret WxPhoneNum
	err := w.wxManager.client.CallWithJson(context.Background(), &ret, "POST", w.wxManager.url("/wxa/business/getuserphonenumber?%v", query.Encode()), nil, codeJson)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
