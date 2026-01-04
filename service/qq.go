package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/global"
	"server/model/other"
	"server/utils"
)

type QQService struct {
}

// GetAccessTokenByCode 通过 Authorization code 获取 Access token
func (service *QQService) GetAccessTokenByCode(code string) (other.AccessTokenResponse, error) {
	data := other.AccessTokenResponse{}
	clientID := global.Config.QQ.AppID
	clientSecret := global.Config.QQ.AppKey
	redirectUri := global.Config.QQ.RedirectURI
	urlStr := "https://graph.qq.com/oauth2.0/token"
	method := "GET"
	params := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
		"redirect_uri":  redirectUri,
		"fmt":           "json",
		"need_openid":   "1",
	}

	// 请求
	resp, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	// 反序列化数据
	if err = json.Unmarshal(byteData, &data); err != nil {
		return other.AccessTokenResponse{}, err
	}
	return data, nil
}

func (service *QQService) GetUserInfoByAccessTokenAndOpenid(accessToken, openid string) (other.UserInfoResponse, error) {
	data := other.UserInfoResponse{}
	oauthConsumerKey := global.Config.QQ.AppID
	urlStr := "https://graph.qq.com/user/get_user_info"
	method := "GET"
	params := map[string]string{
		"access_token":       accessToken,
		"oauth_consumer_key": oauthConsumerKey,
		"openid":             openid,
	}

	resp, err := utils.HttpRequest(urlStr, method, nil, params, nil)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	if err = json.Unmarshal(byteData, &data); err != nil {
		return other.UserInfoResponse{}, err
	}
	return data, nil
}
