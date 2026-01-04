package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"server/global"
	"server/model/request"
	"time"
)

type JWT struct {
	AccessTokenSecret  []byte // access token 的密钥
	RefreshTokenSecret []byte // refresh token 的密钥
}

var (
	TokenExpired     = errors.New("token is expired")           // token 过期
	TokenNotValidYet = errors.New("token not active yet")       // token 还不可用
	TokenMalformed   = errors.New("that's not even a token")    // token 格式错误
	TokenInvalid     = errors.New("couldn't handle this token") // token 无效
)

// NewJWT 创建一个新的 JWT 实例，初始化 AccessToken 和 RefreshToken 密钥
func NewJWT() *JWT {
	return &JWT{
		AccessTokenSecret:  []byte(global.Config.Jwt.AccessTokenSecret),
		RefreshTokenSecret: []byte(global.Config.Jwt.RefreshTokenSecret),
	}
}

// CreateAccessClaims 创建 Access Token 的 Claims，包含基本信息和过期时间等
func (j *JWT) CreateAccessClaims(baseClaims request.BaseClaims) request.JwtCustomClaims {
	// 获取 access token 的过期时间
	ep, _ := ParseDuration(global.Config.Jwt.AccessTokenExpiryTime)

	claims := request.JwtCustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Lumos"}, // 受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	return claims
}

// CreateAccessToken 创建 Access Token
func (j *JWT) CreateAccessToken(claims request.JwtCustomClaims) (string, error) {
	// 创建 jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.AccessTokenSecret)
}

// CreateRefreshClaims 创建 Refresh Token 的 Claims
func (j *JWT) CreateRefreshClaims(baseClaims request.BaseClaims) request.JwtCustomRefreshClaims {
	ep, _ := ParseDuration(global.Config.Jwt.RefreshTokenExpiryTime)

	claims := request.JwtCustomRefreshClaims{
		UserID: baseClaims.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"Lumos"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}
	return claims
}

// CreateRefreshToken 创建 Refresh Token
func (j *JWT) CreateRefreshToken(claims request.JwtCustomRefreshClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

// ParseAccessToken 解析 Access Token，并验证 Claims 信息
func (j *JWT) ParseAccessToken(tokenString string) (*request.JwtCustomClaims, error) {
	// 解析 token
	claims, err := j.parseToken(tokenString, &request.JwtCustomClaims{}, j.AccessTokenSecret)
	if err != nil {
		return nil, err
	}
	// 验证解析出的 Claims 类型是否正确
	if customClaims, ok := claims.(*request.JwtCustomClaims); ok {
		return customClaims, nil
	}
	// 默认返回 TokenInvalid 错误
	return nil, TokenInvalid
}

// ParseRefreshToken 解析 Refresh Token，并验证 Claims 信息
func (j *JWT) ParseRefreshToken(tokenString string) (*request.JwtCustomRefreshClaims, error) {
	claims, err := j.parseToken(tokenString, &request.JwtCustomRefreshClaims{}, j.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}
	if customRefreshClaims, ok := claims.(*request.JwtCustomRefreshClaims); ok {
		return customRefreshClaims, nil
	}
	return nil, TokenInvalid
}

// parseToken 通用的 Token 解析方法，验证 Token 是否有效并返回 Claims
func (j *JWT) parseToken(tokenString string, claims jwt.Claims, secretKey interface{}) (interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) { // 处理 token 的验证错误
			switch {
			case ve.Errors&jwt.ValidationErrorMalformed != 0:
				return nil, TokenMalformed
			case ve.Errors&jwt.ValidationErrorExpired != 0:
				return nil, TokenExpired
			case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
				return nil, TokenNotValidYet
			default:
				return nil, TokenInvalid
			}
		}
		// 默认返回 Token 无效错误
		return nil, TokenInvalid
	}
	// 如果 Token 验证通过，返回 Claims
	if token.Valid {
		return token.Claims, nil
	}
	return nil, TokenInvalid
}
