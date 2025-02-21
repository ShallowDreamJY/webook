package ioc

import "webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service {
	appId := "appdId"
	appSecret := "appdSecret"
	return wechat.NewService(appId, appSecret)
}
