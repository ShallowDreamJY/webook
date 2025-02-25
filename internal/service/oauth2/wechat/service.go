package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webook/internal/domain"
)

var redirectURI, _ = url.PathUnescape("redirectURI")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error)
}

type ServiceImpl struct {
	appId     string
	appSecret string
}

func NewService(appId string, appSecret string) Service {
	return &ServiceImpl{
		appId:     appId,
		appSecret: appSecret,
	}
}

func (s *ServiceImpl) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=%s&scope=snsapi_login&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

func (s *ServiceImpl) VerifyCode(ctx context.Context, code string, state string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	resp, err := http.Get(target)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{},
			fmt.Errorf("微信返回错误信息， 错误码%d， 错误信息%s", res.ErrCode, res.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  res.Openid,
		UnionId: res.UnionId,
	}, nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	Openid  string `json:"openid"`
	UnionId string `json:"unionid"`
	Scope   string `json:"scope"`
}
