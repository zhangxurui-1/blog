package response

import (
	"github.com/gofrs/uuid"
	"server/model/database"
)

// Login 用户登录的返回结构体
type Login struct {
	User                 database.User `json:"user"`
	AccessToken          string        `json:"access_token"`
	AccessTokenExpiresAt int64         `json:"access_token_expires_at"`
}

// UserCard 用户卡片响应
type UserCard struct {
	UUID      uuid.UUID `json:"uuid"`
	UserName  string    `json:"user_name"`
	Avatar    string    `json:"avatar"`
	Address   string    `json:"address"`
	Signature string    `json:"signature"`
}

// UserChart 用户图表回复结构体
type UserChart struct {
	DateList     []string `json:"date_list"`
	LoginData    []int    `json:"login_data"`
	RegisterData []int    `json:"register_data"`
}
