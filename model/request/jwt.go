package request

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"server/model/appTypes"
)

// JwtCustomClaims access token
type JwtCustomClaims struct {
	BaseClaims
	jwt.RegisteredClaims
}

// JwtCustomRefreshClaims refresh token
type JwtCustomRefreshClaims struct {
	UserID uint
	jwt.RegisteredClaims
}

// BaseClaims 结构体用于存储基本的用户信息，作为 JWT 的 Claim 部分
type BaseClaims struct {
	UserID uint
	UUID   uuid.UUID
	RoleID appTypes.RoleID
}
