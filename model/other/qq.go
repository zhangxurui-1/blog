package other

// AccessTokenResponse 通过授权码获取的 Access Token 返回结构
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
}

// UserInfoResponse 用户信息的返回结构
type UserInfoResponse struct {
	Ret          int    `json:"ret"`
	Msg          string `json:"msg"`
	IsLost       bool   `json:"is_lost"`
	Nickname     string `json:"nickname"`
	Figureurl    string `json:"figureurl"`      // 30x30头像URL
	Figureurl1   string `json:"figureurl_1"`    // 50x50头像URL
	Figureurl2   string `json:"figureurl_2"`    // 100x100头像URL
	FigureurlQQ1 string `json:"figureurl_qq_1"` // 40x40 QQ头像URL
	FigureurlQQ2 string `json:"figureurl_qq_2"` // 100x100 QQ头像URL
}
